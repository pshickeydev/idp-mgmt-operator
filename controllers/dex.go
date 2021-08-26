// Copyright Red Hat

package controllers

import (
	"context"
	"fmt"
	"os"

	identitatemdexserverv1alpha1 "github.com/identitatem/dex-operator/api/v1alpha1"
	dexoperatorconfig "github.com/identitatem/dex-operator/config"
	identitatemv1alpha1 "github.com/identitatem/idp-client-api/api/identitatem/v1alpha1"
	"github.com/identitatem/idp-mgmt-operator/deploy"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusteradmapply "open-cluster-management.io/clusteradm/pkg/helpers/apply"
	clusteradmasset "open-cluster-management.io/clusteradm/pkg/helpers/asset"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	dexOperatorImageEnvName string = "DEX_OPERATOR_IMAGE"
)

func (r *AuthRealmReconciler) syncDexCRs(authRealm *identitatemv1alpha1.AuthRealm) error {
	r.Log.Info("syncDexCRs", "AuthRealm.Name", authRealm.Name, "AuthRealm.Namespace", authRealm.Namespace)
	//TODO this test should maybe be in a webhook to avoid the creation of an invalid CR
	if len(authRealm.Spec.IdentityProviders) != 1 {
		return fmt.Errorf("the identityproviders array of the authrealm %s can have only and only one element", authRealm.Name)
	}

	// Create namespace and Install the dex-operator
	if err := r.installDexOperator(authRealm); err != nil {
		return err
	}
	//Create DexServer CR
	if err := r.createDexServer(authRealm); err != nil {
		return err
	}

	return nil
}

func (r *AuthRealmReconciler) installDexOperator(authRealm *identitatemv1alpha1.AuthRealm) error {
	r.Log.Info("installDexOperator", "Name", authRealm.Name, "Namespace", authRealm.Name)

	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.
		WithClient(r.KubeClient, r.APIExtensionClient, r.DynamicClient).
		WithTemplateFuncMap(FuncMap()).
		Build()

	dexOperatorInage := os.Getenv(dexOperatorImageEnvName)
	if len(dexOperatorInage) == 0 {
		return fmt.Errorf("EnvVar %s not provided", dexOperatorImageEnvName)
	}

	//Create the namespace
	readerDeploy := deploy.GetScenarioResourcesReader()
	readerDexOperator := dexoperatorconfig.GetScenarioResourcesReader()

	values := struct {
		Image     string
		AuthRealm *identitatemv1alpha1.AuthRealm
		Reader    *clusteradmasset.ScenarioResourcesReader
		File      string
		NewName   string
	}{
		Image:     dexOperatorInage,
		AuthRealm: authRealm,
		Reader:    readerDexOperator,
		File:      "rbac/role.yaml",
		NewName:   "dex-operator-manager-role",
	}

	files := []string{
		"dex-operator/namespace.yaml",
		"dex-operator/service_account.yaml",
		"dex-operator/role.yaml",
		"dex-operator/role_binding.yaml",
	}

	_, err := applier.ApplyDirectly(readerDeploy, values, false, "", files...)
	if err != nil {
		return err
	}

	_, err = applier.ApplyDeployments(readerDeploy, values, false, "", "dex-operator/manager.yaml")
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRealmReconciler) installDexCRDs() error {
	r.Log.Info("installDexCRDs")

	applierBuilder := &clusteradmapply.ApplierBuilder{}
	applier := applierBuilder.
		WithClient(r.KubeClient, r.APIExtensionClient, r.DynamicClient).
		WithTemplateFuncMap(FuncMap()).
		Build()

	readerIDPStrategyOperator := dexoperatorconfig.GetScenarioResourcesReader()
	files := []string{"crd/bases/auth.identitatem.io_dexclients.yaml",
		"crd/bases/auth.identitatem.io_dexservers.yaml"}

	_, err := applier.ApplyDirectly(readerIDPStrategyOperator, nil, false, "", files...)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRealmReconciler) createDexServer(authRealm *identitatemv1alpha1.AuthRealm) error {
	r.Log.Info("createDexServer", "Name", authRealm.Name, "Namespace", authRealm.Name)
	dexServerExists := true
	dexServer := &identitatemdexserverv1alpha1.DexServer{}
	if err := r.Client.Get(context.TODO(), client.ObjectKey{Name: authRealm.Name, Namespace: authRealm.Name}, dexServer); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		dexServerExists = false
		dexServer = &identitatemdexserverv1alpha1.DexServer{
			ObjectMeta: metav1.ObjectMeta{
				Name:      authRealm.Name,
				Namespace: authRealm.Name,
			},
		}
	}

	if err := r.updateDexServer(authRealm, dexServer); err != nil {
		return err
	}

	switch dexServerExists {
	case true:
		r.Log.V(1).Info("createDexServer update dexServer", "Name", authRealm.Name, "Namespace", authRealm.Name)
		if err := r.Client.Update(context.TODO(), dexServer); err != nil {
			return err
		}
	case false:
		r.Log.V(1).Info("createDexServer create dexServer", "Name", authRealm.Name, "Namespace", authRealm.Name)
		if err := r.Client.Create(context.TODO(), dexServer); err != nil {
			return err
		}
	}

	return nil
}

func (r *AuthRealmReconciler) updateDexServer(authRealm *identitatemv1alpha1.AuthRealm, dexServer *identitatemdexserverv1alpha1.DexServer) error {
	r.Log.Info("updateDexServer", "Name", authRealm.Name, "Namespace", authRealm.Name)
	dexServer.Spec.Issuer = authRealm.Spec.Host
	if len(authRealm.Spec.CertificatesSecretRef.Name) != 0 {
		certSecret := &corev1.Secret{}
		if err := r.Client.Get(context.TODO(),
			client.ObjectKey{Name: authRealm.Spec.CertificatesSecretRef.Name, Namespace: authRealm.Namespace},
			certSecret); err != nil {
			return err
		}
		dexServer.Spec.Web.TlsCert = string(certSecret.Data["tls.crt"])
		dexServer.Spec.Web.TlsKey = string(certSecret.Data["tls.key"])
	}
	cs, err := r.createDexConnectors(authRealm)
	if err != nil {
		return err
	}
	dexServer.Spec.Connectors = cs
	return nil
}

func (r *AuthRealmReconciler) createDexConnectors(authRealm *identitatemv1alpha1.AuthRealm) (cs []identitatemdexserverv1alpha1.ConnectorSpec, err error) {
	r.Log.Info("createDexConnectors", "Name", authRealm.Name, "Namespace", authRealm.Name)

	for _, idp := range authRealm.Spec.IdentityProviders {
		cs = make([]identitatemdexserverv1alpha1.ConnectorSpec, 0)
		if idp.GitHub != nil {
			c, err := r.createConnector(authRealm, "github", idp.GitHub.ClientSecret.Name)
			if err != nil {
				return nil, err
			}
			cs = append(cs, *c)
		}
		if idp.Google != nil {
			c, err := r.createConnector(authRealm, "google", idp.Google.ClientSecret.Name)
			if err != nil {
				return nil, err
			}
			cs = append(cs, *c)
		}
		if idp.HTPasswd != nil {
			c, err := r.createConnector(authRealm, "htppasswd", "")
			if err != nil {
				return nil, err
			}
			cs = append(cs, *c)
		}
		if idp.LDAP != nil {
			c, err := r.createConnector(authRealm, "ldap", idp.LDAP.BindPassword.Name)
			if err != nil {
				return nil, err
			}
			cs = append(cs, *c)
		}
		if idp.OpenID != nil {
			c, err := r.createConnector(authRealm, "oidc", idp.OpenID.ClientSecret.Name)
			if err != nil {
				return nil, err
			}
			cs = append(cs, *c)
		}
		if len(cs) == 0 {
			return nil, fmt.Errorf("no identityProvider defined in %s/%s", authRealm.Name, authRealm.Name)
		}
	}
	return cs, err
}

func (r *AuthRealmReconciler) createConnector(authRealm *identitatemv1alpha1.AuthRealm,
	identityProviderType string, clientSecretName string) (c *identitatemdexserverv1alpha1.ConnectorSpec, err error) {

	c = &identitatemdexserverv1alpha1.ConnectorSpec{
		Type: identityProviderType,
		Name: authRealm.Name,
		Id:   authRealm.Name,
		Config: identitatemdexserverv1alpha1.ConfigSpec{
			ClientID:        authRealm.Name,
			ClientSecretRef: clientSecretName,
		},
	}

	return c, nil
}

func (r *AuthRealmReconciler) deleteAuthRealmNamespace(authRealm *identitatemv1alpha1.AuthRealm) error {
	ns := &corev1.Namespace{}
	if err := r.Client.Get(context.TODO(), client.ObjectKey{Name: authRealm.Name}, ns); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return r.Client.Delete(context.TODO(), ns)
}