package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v13 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func RHSSODeployment(cr *v1alpha1.Keycloak) *v13.StatefulSet {
	return &v13.StatefulSet{
		ObjectMeta: v12.ObjectMeta{
			Name:      KeycloakDeploymentName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":       ApplicationName,
				"component": KeycloakDeploymentComponent,
			},
		},
		Spec: v13.StatefulSetSpec{
			Selector: &v12.LabelSelector{
				MatchLabels: map[string]string{
					"app":       ApplicationName,
					"component": KeycloakDeploymentComponent,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: v12.ObjectMeta{
					Name:      KeycloakDeploymentName,
					Namespace: cr.Namespace,
					Labels: map[string]string{
						"app":       ApplicationName,
						"component": KeycloakDeploymentComponent,
					},
				},
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: ServingCertSecretName,
							VolumeSource: v1.VolumeSource{
								Secret: &v1.SecretVolumeSource{
									SecretName: ServingCertSecretName,
									Optional:   &[]bool{true}[0],
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:  KeycloakDeploymentName,
							Image: RHSSOImage,
							Ports: []v1.ContainerPort{
								{
									ContainerPort: KeycloakServicePort,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 8080,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9990,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 8778,
									Protocol:      "TCP",
								},
							},
							Env: []v1.EnvVar{
								// Database settings
								{
									Name:  "DB_SERVICE_PREFIX_MAPPING",
									Value: PostgresqlServiceName + "=DB",
								},
								{
									Name:  "TX_DATABASE_PREFIX_MAPPING",
									Value: PostgresqlServiceName + "=DB",
								},
								{
									Name:  "DB_JNDI",
									Value: "java:jboss/datasources/KeycloakDS",
								},
								{
									Name:  "DB_SCHEMA",
									Value: "public",
								},
								{
									Name: "DB_USERNAME",
									ValueFrom: &v1.EnvVarSource{
										SecretKeyRef: &v1.SecretKeySelector{
											LocalObjectReference: v1.LocalObjectReference{
												Name: DatabaseSecretName,
											},
											Key: DatabaseSecretUsernameProperty,
										},
									},
								},
								{
									Name: "DB_PASSWORD",
									ValueFrom: &v1.EnvVarSource{
										SecretKeyRef: &v1.SecretKeySelector{
											LocalObjectReference: v1.LocalObjectReference{
												Name: DatabaseSecretName,
											},
											Key: DatabaseSecretPasswordProperty,
										},
									},
								},
								{
									Name:  "DB_DATABASE",
									Value: PostgresqlDatabase,
								},
								// Discovery settings
								{
									Name:  "JGROUPS_PING_PROTOCOL",
									Value: "dns.DNS_PING",
								},
								{
									Name:  "OPENSHIFT_DNS_PING_SERVICE_NAME",
									Value: KeycloakDiscoveryServiceName + "." + cr.Namespace + ".svc.cluster.local",
								},
								{
									Name: "SSO_ADMIN_USERNAME",
									ValueFrom: &v1.EnvVarSource{
										SecretKeyRef: &v1.SecretKeySelector{
											LocalObjectReference: v1.LocalObjectReference{
												Name: "credential-" + cr.Name,
											},
											Key: AdminUsernameProperty,
										},
									},
								},
								{
									Name: "SSO_ADMIN_PASSWORD",
									ValueFrom: &v1.EnvVarSource{
										SecretKeyRef: &v1.SecretKeySelector{
											LocalObjectReference: v1.LocalObjectReference{
												Name: "credential-" + cr.Name,
											},
											Key: AdminPasswordProperty,
										},
									},
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      ServingCertSecretName,
									MountPath: "/etc/x509/https",
								},
							},
							LivenessProbe: &v1.Probe{
								InitialDelaySeconds: 60,
								TimeoutSeconds:      1,
								Handler: v1.Handler{
									HTTPGet: &v1.HTTPGetAction{
										Path:   "/auth/realms/master",
										Port:   intstr.FromInt(8080),
										Scheme: "HTTP",
									},
								},
							},
							ReadinessProbe: &v1.Probe{
								TimeoutSeconds:      1,
								InitialDelaySeconds: 10,
								Handler: v1.Handler{
									HTTPGet: &v1.HTTPGetAction{
										Path:   "/auth/realms/master",
										Port:   intstr.FromInt(8080),
										Scheme: "HTTP",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func RHSSODeploymentSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      KeycloakDeploymentName,
		Namespace: cr.Namespace,
	}
}

func RHSSODeploymentReconciled(cr *v1alpha1.Keycloak, currentState *v13.StatefulSet) *v13.StatefulSet {
	reconciled := currentState.DeepCopy()
	reconciled.ResourceVersion = currentState.ResourceVersion
	reconciled.Spec.Template.Spec.Volumes = []v1.Volume{
		{
			Name: ServingCertSecretName,
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: ServingCertSecretName,
					Optional:   &[]bool{true}[0],
				},
			},
		},
	}
	reconciled.Spec.Template.Spec.Containers = []v1.Container{
		{
			Name:  KeycloakDeploymentName,
			Image: RHSSOImage,
			Ports: []v1.ContainerPort{
				{
					ContainerPort: KeycloakServicePort,
					Protocol:      "TCP",
				},
				{
					ContainerPort: 8080,
					Protocol:      "TCP",
				},
				{
					ContainerPort: 9990,
					Protocol:      "TCP",
				},
				{
					ContainerPort: 8778,
					Protocol:      "TCP",
				},
			},
			VolumeMounts: []v1.VolumeMount{
				{
					Name:      ServingCertSecretName,
					MountPath: "/etc/x509/https",
				},
			},
			LivenessProbe: &v1.Probe{
				InitialDelaySeconds: 60,
				TimeoutSeconds:      1,
				Handler: v1.Handler{
					HTTPGet: &v1.HTTPGetAction{
						Path:   "/auth/realms/master",
						Port:   intstr.FromInt(8080),
						Scheme: "HTTP",
					},
				},
			},
			ReadinessProbe: &v1.Probe{
				TimeoutSeconds:      1,
				InitialDelaySeconds: 10,
				Handler: v1.Handler{
					HTTPGet: &v1.HTTPGetAction{
						Path:   "/auth/realms/master",
						Port:   intstr.FromInt(8080),
						Scheme: "HTTP",
					},
				},
			},
			Env: []v1.EnvVar{
				// Database settings
				{
					Name:  "DB_SERVICE_PREFIX_MAPPING",
					Value: PostgresqlServiceName + "=DB",
				},
				{
					Name:  "TX_DATABASE_PREFIX_MAPPING",
					Value: PostgresqlServiceName + "=DB",
				},
				{
					Name:  "DB_JNDI",
					Value: "java:jboss/datasources/KeycloakDS",
				},
				{
					Name:  "DB_SCHEMA",
					Value: "public",
				},
				{
					Name: "DB_USERNAME",
					ValueFrom: &v1.EnvVarSource{
						SecretKeyRef: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: DatabaseSecretName,
							},
							Key: DatabaseSecretUsernameProperty,
						},
					},
				},
				{
					Name: "DB_PASSWORD",
					ValueFrom: &v1.EnvVarSource{
						SecretKeyRef: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: DatabaseSecretName,
							},
							Key: DatabaseSecretPasswordProperty,
						},
					},
				},
				{
					Name:  "DB_DATABASE",
					Value: PostgresqlDatabase,
				},
				// Discovery settings
				{
					Name:  "JGROUPS_PING_PROTOCOL",
					Value: "dns.DNS_PING",
				},
				{
					Name:  "OPENSHIFT_DNS_PING_SERVICE_NAME",
					Value: KeycloakDiscoveryServiceName + "." + cr.Namespace + ".svc.cluster.local",
				},
				{
					Name: "SSO_ADMIN_USERNAME",
					ValueFrom: &v1.EnvVarSource{
						SecretKeyRef: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "credential-" + cr.Name,
							},
							Key: AdminUsernameProperty,
						},
					},
				},
				{
					Name: "SSO_ADMIN_PASSWORD",
					ValueFrom: &v1.EnvVarSource{
						SecretKeyRef: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "credential-" + cr.Name,
							},
							Key: AdminPasswordProperty,
						},
					},
				},
			},
		},
	}

	return reconciled
}
