# SVC Bridge

Forward messages published on one URI to another, maybe on another namespace.

The params file needs to specify `fromuri` and `touri` e.g:

```
fromuri: q_IwWMaqz1xVDiAUhwsrRNSTa9dToiS1acxcAQ94G2U=/sensors/mpa
touri: mask/sensors
```

Due to laziness by the developer, do not use a namespace alias in the fromuri.
