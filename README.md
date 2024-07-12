<img alt="CyberArk Banner" src="images/cyberark-banner.jpg">

# CyberArk Privilege Access Manager - Go SDK

<!--
Author:   David Hisel <david.hisel@cyberark.com>
Updated:  <2024-07-12 13:45:25 david.hisel>
-->

## Summary

This Go SDK interfaces with CyberArk Privilege Cloud REST endpoints.

Here is [a link to the documentation](https://docs.cyberark.com/privilege-cloud-shared-services/latest/en/Content/WebServices/Implementing%20Privileged%20Account%20Security%20Web%20Services%20.htm) for reference.

Note that this SDK is a **WORK IN PROGRESS**, and does not implement all the endpoints, yet.

## Examples

Look in [`example/`](./example) folder for reference implementations.

To run an example, follow these steps:

1. Copy the `creds.toml.example` to `creds.toml`
2. Edit `creds.toml` and add your values
3. Run the example from the main project directory

    ```shell
    go run example/session/main.go
    ```

## License

Copyright (c) 2024 CyberArk Software Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

<http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

For the full license text see [`LICENSE`](LICENSE).

## Contributing

We welcome contributions of all kinds to this repository. For
instructions on how to get started and descriptions of our development
workflows, please see our [contributing
guide](CONTRIBUTING.md).

[Code of Conduct](CODE_OF_CONDUCT.md).
