# `diamcircle/expingest-verify-range`

This docker image allows running multiple instances of `aurora ingest verify-command` on a single machine or running it in [AWS Batch](https://aws.amazon.com/batch/).

## Env variables

### Running locally

| Name     | Description                                           |
|----------|-------------------------------------------------------|
| `BRANCH` | Git branch to build (useful for testing PRs)          |
| `FROM`   | First ledger of the range (must be checkpoint ledger) |
| `TO`     | Last ledger of the range (must be checkpoint ledger)  |

### Running in AWS Batch

| Name                 | Description                                                          |
|----------------------|----------------------------------------------------------------------|
| `BRANCH`             | Git branch to build (useful for testing PRs)                         |
| `BATCH_START_LEDGER` | First ledger of the AWS Batch Job, must be a checkpoint ledger or 1. |
| `BATCH_SIZE`         | Size of the batch, must be multiple of 64.                           |

#### Example

When you start 10 jobs with `BATCH_START_LEDGER=63` and `BATCH_SIZE=64`
it will run the following ranges:

| `AWS_BATCH_JOB_ARRAY_INDEX` | `FROM` | `TO` |
|-----------------------------|--------|------|
| 0                           | 63     | 127  |
| 1                           | 127    | 191  |
| 2                           | 191    | 255  |
| 3                           | 255    | 319  |

## Tips when using AWS Batch

* In "Job definition" set vCPUs to 2 and Memory to 4096. This represents the `c5.large` instances Aurora should be using.
* In "Compute environments":
    * Set instance type to "c5.large".
    * Set "Maximum vCPUs" to 2x the number of instances you want to start (because "c5.large" has 2 vCPUs). Ex. 10 vCPUs = 5 x "c5.large" instances.
* Use spot instances! It's much cheaper and speed of testing will be the same in 99% of cases.
* You need to publish the image if there are any changes in `Dockerfile` or one of the scripts.
* When batch processing is over check if instances have been terminated. Sometimes AWS doesn't terminate them.
* Make sure the job timeout is set to a larger value if you verify larger ranges. Default is just 100 seconds.
