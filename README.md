# Change value in file

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-change-value?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-change-value/releases)

Changes a selected value in a targeted file.

<details>
<summary>Description</summary>

This Step changes a selected value in a targeted file, for example, constants.

### Configuring the Step

1. Set the **File path** in which you want to change a value.
2. Set the **Current value** which needs to be changed.
3. Set the **New value**.

### Troubleshooting
Please make sure that your targeted file path is correct, existing and relative from the root folder.
Please make sure that you set a correct **Current value** which exists in the targeted file. You must define this value.
Please make sure that you define the New value.

### Related Steps
- [Change Working Directory for subsequent Steps](https://www.bitrise.io/integrations/steps/change-workdir)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `file` | Path to the file in which the value replacement will be performed. The path must be relative to the repository root. | required |  |
| `old_value` | The exact string to search for and replace in the target file. The step fails if this value is not found and Mark the Step as Failed if not found is enabled. | required |  |
| `new_value` | The string that replaces all occurrences of the current value in the target file. | required |  |
| `show_file` | When enabled, prints the full file content to the build log before and after the replacement. |  | `false` |
| `notfound_exit` | When enabled, the step fails if the current value is not found in the target file. Disable this to allow the step to succeed even if no replacement was made. |  | `true` |
| `use_sudo` | When enabled, the step retries reading and writing the target file with `sudo` if it hits a permission error.  This is required on stacks where the step runs as the non-root user and the target file is owned by `root`. Disable this if `sudo` is unavailable in your environment or you want permission errors to fail the step directly. |  | `true` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-change-value/pulls) and [issues](https://github.com/bitrise-steplib/steps-change-value/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
