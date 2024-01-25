// Package txsub provides the machinery that aurora uses to submit transactions to
// the diamcircle network and track their progress.  It also helps to hide some of the
// complex asynchronous nature of transaction submission, waiting to respond to
// submitters when no definitive state is known.
package txsub

// Package layout:
// - main.go: interface and result types
// - errors.go: error definitions exposed by txsub
// - system.go: txsub.System, the struct that ties all the interfaces together
// - open_submission_list.go: A default implementation of the OpenSubmissionList interface
// - submitter.go: A default implementation of the Submitter interface
