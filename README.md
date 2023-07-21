# Aegis Helm Chart

Aegis keeps your secrets secret. With Aegis, you can rest assured that your sensitive data is always secure and protected. Aegis is perfect for securely storing arbitrary configuration information at a central location and securely dispatching it to workloads.

## Installation

To use Aegis, follow the steps below:

1. Add Aegis Helm repository:

    ```bash
    helm repo add aegis https://abhishek44sharma.github.io/aegis/
    ```

2. Install Aegis using Helm:

    ```bash
    helm install aegis aegis/helm-charts --version 0.1.0
    ```

## Options

The following options can be passed to the `helm install` command to set global variables:

- `--set deploySpire=<true/false>`: This flag can be passed to install or skip Spire.
- `--set platform=<istanbul/istanbul-fips/photon/photos-fips>`: This flag can be passed to install Aegis with the given platform Docker image.

Here's an example command with the above options:

```bash
helm install aegis aegis/helm-charts --version 0.1.0 --set deploySpire=true --set platform=istanbul
```

Make sure to replace `<true/false>` and `<istanbul/istanbul-fips/photon/photos-fips>` with the desired values.

## License

This project is licensed under the [MIT License](LICENSE).