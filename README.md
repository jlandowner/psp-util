# Pod Security Policy Utility
A Kubectl extention utility to manage `Pod Security Policy(PSP)` and RBAC Resources.

Attach/Detach PSP to RBACs(Group, User) or ServiceAccounts and view the relations which PSP is effected to each Subjects.

See the details of PSP: 
- https://kubernetes.io/docs/concepts/policy/pod-security-policy/

See the Best Practices of PSP: 
- https://aws.github.io/aws-eks-best-practices/pods/#recommendations
- https://github.com/sysdiglabs/kube-psp-advisor

# Install

You can install it by [krew](https://krew.sigs.k8s.io).
After [installing krew](https://krew.sigs.k8s.io/docs/user-guide/setup/install/), run the following command:

```shell
kubectl krew install --manifest-url https://raw.githubusercontent.com/jlandowner/psp-util/master/psp-util.yaml

```

# Usage

```shell
$ kubectl psp-util

A Kubectl extention utility to manage Pod Security Policy(PSP) and related RBAC Resources.
Attach/Detach PSP to/from RBACs(Group, User) or ServiceAccounts and
view the relations which PSP is effected to the Subjects in cluster.

Complete documentation is available at http://github.com/jlandowner/psp-util

Usage:
  psp-util [command]

Available Commands:
  attach      Attach PSP to RBAC Subject (Auto generate managed ClusterRole and ClusterRoleBinding)
  clean       Clean managed ClusterRole and ClusterRoleBinding
  detach      Detach PSP from RBAC Subject
  help        Help about any command
  list        List PSP and RBAC associated with it.
  tree        View relational tree between PSP and Subjects
  version     Print the version number

Flags:
  -h, --help                help for psp-util
      --kubeconfig string   kube config file (default is $HOME/.kube/config)

Use "psp-util [command] --help" for more information about a command.
```

# Command details
## list

`list` shows all PSPs in cluster, and also ClusterRoles and ClusterRoleBindings associated with each of them.

A column `PSP-UTIL-MANAGED` is whether these ClusterRoles and ClusterRoleBindings are auto-created and managed by `psp-util`.

```shell
$ kubectl psp-util list
PSP-NAME                                 CLUSTER-ROLE                                      CLUSTER-ROLE-BINDING                              PSP-UTIL-MANAGED
eks.privileged                           eks:podsecuritypolicy:privileged                  eks:podsecuritypolicy:authenticated               false
pod-security-policy-all-20200702180710   psp-util.pod-security-policy-all-20200702180710   psp-util.pod-security-policy-all-20200702180710   true
restricted                               psp-util.restricted                               psp-util.restricted                               true

```

## tree

`tree` shows the relations between PSP and Subjects by tree expressions.

```shell
$ kubectl psp-util tree
📙 PSP eks.privileged
└── 📕 ClusterRole eks:podsecuritypolicy:privileged
    └── 📘 ClusterRoleBinding eks:podsecuritypolicy:authenticated
        └── 📗 Subject{Kind: Group, Name: system:master, Namespace: }
        └── 📗 Subject{Kind: ServiceAccount, Name: default, Namespace: kube-system}

📙 PSP pod-security-policy-all-20200702180710
└── 📕 ClusterRole psp-util.pod-security-policy-all-20200702180710
    └── 📘 ClusterRoleBinding psp-util.pod-security-policy-all-20200702180710
        └── 📗 Subject{Kind: Group, Name: system:authenticated, Namespace: }

📙 PSP restricted
└── 📕 ClusterRole psp-util.restricted
    └── 📘 ClusterRoleBinding psp-util.restricted
        └── 📗 Subject{Kind: Group, Name: my:group, Namespace: }
        └── 📗 Subject{Kind: ServiceAccount, Name: default, Namespace: default}

```

## attach

`attach` attaches PSP to Subjects(Group, User or ServiceAccount).

```shell
Usage:
  psp-util attach PSP-NAME [ --group | --user | --sa ] SUBJECT-NAME [flags]

Flags:
  -g, --group string       set Subject's Name and use Kind Group
  -u, --user string        set Subject's Name and use Kind User
  -s, --sa string          set Subject's Name and use Kind ServiceAccount
  -n, --namespace string   set Subject's Namespace (only used when kind is ServiceAccount)
      --api-group string   set Subject's APIGroup
      --kind string        set Subject's Kind
      --name string        set Subject's Name
```

If there is no managed ClusterRole and ClusterRoleBinding associated with the given PSP, 
it will generate them automaticaly.

### Examples

Attach MyPSP to Subject{APIGroup: rbac.authorization.k8s.io, Kind: Group, Name: system:authenticated}.

```shell
$ kubectl psp-util attach MyPSP --group system:authenticated
Managed ClusterRole is not found...Created
Managed ClusterRoleBinding is not found...Created

```

Attach MyPSP to Subject{Kind: ServiceAccount, Name: default, Namespace: kube-system}.

```shell
$ kubectl psp-util attach MyPSP --sa default -n kube-system
```

Or, set all subject's info directly.

```shell
$ kubectl psp-util attach MyPSP --api-group=rbac.authorization.k8s.io --kind=Group --name=system:authenticated
```


## detach

`detach` detached a Subject from PSP.

It removes the Subject from the ClusterRoleBinding only if there is a managed ClusterRoleBinding in cluster.

All the options are the same as for the `attach` command.

```shell
Usage:
  psp-util detach PSP-NAME [ --group | --user | --sa ] SUBJECT-NAME [flags]

Flags:
  -g, --group string       set Subject's Name and use Kind Group
  -u, --user string        set Subject's Name and use Kind User
  -s, --sa string          set Subject's Name and use Kind ServiceAccount
  -n, --namespace string   set Subject's Namespace (only used when kind is ServiceAccount)
      --api-group string   set Subject's APIGroup
      --kind string        set Subject's Kind
      --name string        set Subject's Name
```

## clean

`clean` delete a managed ClusterRole and ClusterRoleBinding.

```shell
Usage:
  psp-util clean PSP-NAME
```

# LICENSE
Apache License Version 2.0 Copyright 2020 jlandowner