# Testing New Ingestion System

 This document describes what you need to know to start testing Aurora's new ingestion system. This system will soon be standard, and your feedback is valuable. If you experience a problem not covered in this document, please add an issue in this repository. For questions please use [Stack Exchange](https://diamcircle.stackexchange.com) or ask it in one of our online [communities](https://www.diamcircle.org/community/#communities).

Please remember that this is still an experimental version of the system and should only be used in staging and test environments!

Thank you for helping us test the new ingestion system in Aurora.

## What is the new ingestion system?

The new ingestion system solves issues found in the previous version like: inconsistent data, relying on Diamcircle-Core database directly, slow responses for specific queries, etc. It allows new kind of features in Aurora, ex. faster path-finding. We published a [blog post](https://www.diamcircle.org/developers-blog/our-new-aurora-ingestion-engine) with more details, please check it out!

## Why would you want to upgrade?

* Ingestion can now run on multiple servers, which means that even if one of your ingesting instances is down, the ingestion will continue on other instances.
* New features like faster path-finding, accounts for signer endpoint, finding all accounts with a given asset, etc. More new features (and plugins!) to come.
* Ingestion does not generate a high load on Diamcircle-Core database.
* With batched requests (not implemented yet) you can get a consistent snapshot of the latest ledger data. Previously this wasn't possible, because some entries were loaded from Diamcircle-Core database and others from the Aurora database, and these could be at different ledgers.
* We will continue to update Aurora 0.* releases with security fixes until end-of-life, but the 1.x release will become the default and recommended version soon. It's better to test this now within your organization. And again, use this release in staging environments only!

## Before you upgrade

* You can rollback to the older version but only when using alpha or beta versions. We won't support rolling back when a stable version is released. To rollback: migrate DB down, rollback to the previous version and run `aurora db init-asset-stats` to regenerate assets stats in your DB.
* If you were using the new ingestion in one of the previous versions of Aurora, you must first remove `ENABLE_EXPERIMENTAL_INGESTION` feature flag and restart all Aurora instances before deploying a new version.
* The init stage (state ingestion) for the public Diamcircle network requires around 1.5GB of RAM. This memory is released after the state ingestion. State ingestion is performed only once. Restarting the server will not trigger it unless Aurora has been upgraded to a newer version (with an updated ingestion pipeline). We are evaluating alternative solutions to making these RAM requirements smaller.It's worth noting that the required memory will become smaller and smaller as more of the buckets in the history archive become [CAP-20](https://github.com/diamcircle/diamcircle-protocol/blob/master/core/cap-0020.md) compatible.
* The CPU footprint of the new ingestion is modest. We were able to successfully run ingestion on an [AWS `c5.large`](https://aws.amazon.com/ec2/instance-types/c5/) instance. However, we highly recommend `c5.xlarge` instances. The init stage takes a few minutes on `c5.xlarge`. `c5.xlarge` is the equivalent of 4 vCPUs and 8GB of RAM. The definition of vCPU for the c5 large family in AWS is the following:
> The 2nd generation Intel Xeon Scalable Processors (Cascade Lake) or 1st generation Intel Xeon Platinum 8000 series (Skylake-SP) processor with a sustained all core Turbo frequency of up to 3.4GHz, and single core turbo frequency of up to 3.5 GHz.

* The state data requires an additional 6GB DB disk space for the public Diamcircle network (as of January 2020). The disk usage will increase when the number of Diamcircle ledger entries increases.
  * `accounts_signers` table: 2340 MB
  * `trust_lines` table: 2052 MB
  * `accounts` table: 1545 MB
  * `offers` table: 61 MB
  * `accounts_data` table: 15 MB
  * `exp_asset_stats` table: less than 1 MB
* A new environment variable (or command line flag) needs to be set so that Aurora can ingest state from the history archives:
   * `HISTORY_ARCHIVE_URLS="archive1,archive2,archive3"` (if you don't have your own pubnet history archive, you can use one of SDF's archives, for example `https://history.diamcircle.org/prd/core-live/core_live_001`)
* Aurora serves the endpoints `/paths` and `/order_book` from an in-memory graph, which is only available on ingesting instances. If some of the instances in your cluster are not configured to ingest, you can configure your proxy server to route those endpoints to the ingesting instances. This is beyond the scope of this document - consult the relevant documentation for your proxy server.

 ## Troubleshooting

### Some endpoints are not available during state ingestion

Endpoints that display state information are not available during initial state ingestion and will return a `503 Service Unavailable`/`Still Ingesting` error.  An example is the `/paths` endpoint (built using offers). Such endpoints will become available after state ingestion is done (usually within a couple of minutes).

### State ingestion is taking a lot of time

State ingestion shouldn't take more than a couple of minutes on an AWS `c5.xlarge` instance, or equivalent.

 It's possible that the progress logs (see below) will not show anything new for a longer period of time or print a lot of progress entries every few seconds. This happens because of the way history archives are designed. The ingestion is still working but it's processing entries of type `DEADENTRY`'. If there is a lot of them in the bucket, there are no _active_ entries to process. We plan to improve the progress logs to display actual percentage progress so it's easier to estimate ETA.

If you see that ingestion is not proceeding for a very long period of time:
1. Check the RAM usage on the machine. It's possible that system run out of RAM and it using swap memory that is extremely slow.
2. If above is not the case, file a new issue in this repository.

### CPU usage goes high every few minutes

This is _by design_. Aurora runs a state verifier routine that compares state in local storage to history archives every 64 ledgers to ensure data changes are applied correctly. If data corruption is detected Aurora will block access to endpoints serving invalid data.

We recommend to keep this security feature turned on however if it's causing problems (due to CPU usage) this can be disabled by `--ingest-disable-state-verification` CLI param or `INGEST-DISABLE-STATE-VERIFICATION` env variable.

### I see `Waiting for the next checkpoint...` messages

If you were running the new system in the past (`ENABLE_EXPERIMENTAL_INGESTION` flag) it's possible that the old and new systems are not in sync. In such case, the upgrade code will activate and will make sure the data is in sync. When this happens you may see `Waiting for the next checkpoint...` messages for up to 5 minutes.

## Reading the logs

In order to check the progress and the status of ingestion you should check the logs. All logs connected to experimental ingestion are tagged with `service=ingest`.

It starts with informing you about state ingestion:
```
INFO[2019-08-29T13:04:13.473+02:00] Starting ingestion system from empty state...  pid=5965 service=ingest temp_set="*io.MemoryTempSet"
INFO[2019-08-29T13:04:15.263+02:00] Reading from History Archive Snapshot         ledger=25565887 pid=5965 service=ingest
```
During state ingestion, Aurora will log number of processed entries every 100,000 entries (there are currently around 7M entries in the public network):
```
INFO[2019-08-29T13:04:34.652+02:00] Processing entries from History Archive Snapshot  ledger=25565887 numEntries=100000 pid=5965 service=ingest
INFO[2019-08-29T13:04:38.487+02:00] Processing entries from History Archive Snapshot  ledger=25565887 numEntries=200000 pid=5965 service=ingest
INFO[2019-08-29T13:04:41.322+02:00] Processing entries from History Archive Snapshot  ledger=25565887 numEntries=300000 pid=5965 service=ingest
INFO[2019-08-29T13:04:48.429+02:00] Processing entries from History Archive Snapshot  ledger=25565887 numEntries=400000 pid=5965 service=ingest
INFO[2019-08-29T13:05:00.306+02:00] Processing entries from History Archive Snapshot  ledger=25565887 numEntries=500000 pid=5965 service=ingest
```
When state ingestion is finished it will proceed to ledger ingestion starting from the next ledger after checkpoint ledger (25565887+1 in this example) to update the state using transaction meta:
```
INFO[2019-08-29T13:39:41.590+02:00] Processing entries from History Archive Snapshot  ledger=25565887 numEntries=5300000 pid=5965 service=ingest
INFO[2019-08-29T13:39:44.518+02:00] Processing entries from History Archive Snapshot  ledger=25565887 numEntries=5400000 pid=5965 service=ingest
INFO[2019-08-29T13:39:47.488+02:00] Processing entries from History Archive Snapshot  ledger=25565887 numEntries=5500000 pid=5965 service=ingest
INFO[2019-08-29T13:40:00.670+02:00] Processed ledger                              ledger=25565887 pid=5965 service=ingest type=state_pipeline
INFO[2019-08-29T13:40:00.670+02:00] Finished processing History Archive Snapshot  duration=2145.337575904 ledger=25565887 numEntries=5529931 pid=5965 service=ingest shutdown=false
INFO[2019-08-29T13:40:00.693+02:00] Reading new ledger                            ledger=25565888 pid=5965 service=ingest
INFO[2019-08-29T13:40:00.694+02:00] Processing ledger                             ledger=25565888 pid=5965 service=ingest type=ledger_pipeline updating_database=true
INFO[2019-08-29T13:40:00.779+02:00] Processed ledger                              ledger=25565888 pid=5965 service=ingest type=ledger_pipeline
INFO[2019-08-29T13:40:00.779+02:00] Finished processing ledger                    duration=0.086024492 ledger=25565888 pid=5965 service=ingest shutdown=false transactions=14
INFO[2019-08-29T13:40:00.815+02:00] Reading new ledger                            ledger=25565889 pid=5965 service=ingest
INFO[2019-08-29T13:40:00.816+02:00] Processing ledger                             ledger=25565889 pid=5965 service=ingest type=ledger_pipeline updating_database=true
INFO[2019-08-29T13:40:00.881+02:00] Processed ledger                              ledger=25565889 pid=5965 service=ingest type=ledger_pipeline
INFO[2019-08-29T13:40:00.881+02:00] Finished processing ledger                    duration=0.06619956 ledger=25565889 pid=5965 service=ingest shutdown=false transactions=29
INFO[2019-08-29T13:40:00.901+02:00] Reading new ledger                            ledger=25565890 pid=5965 service=ingest
INFO[2019-08-29T13:40:00.902+02:00] Processing ledger                             ledger=25565890 pid=5965 service=ingest type=ledger_pipeline updating_database=true
INFO[2019-08-29T13:40:00.972+02:00] Processed ledger                              ledger=25565890 pid=5965 service=ingest type=ledger_pipeline
INFO[2019-08-29T13:40:00.972+02:00] Finished processing ledger                    duration=0.071039012 ledger=25565890 pid=5965 service=ingest shutdown=false transactions=20
```
