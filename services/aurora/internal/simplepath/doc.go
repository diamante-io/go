// Package simplepath provides an implementation of paths. Finder that performs
// a breadth first search for paths against an orderbook.
//
// The core algorithm works as follows:
// 1. `search` object contains a queue of currently extended paths. Queue is
//    initialized with a single-asset path containing destination asset. Paths
//    are linked lists, head of the path is a newly extended source asset, tail
//    of the path is always the destination asset.
// 2. Every iteration of search does the following:
//    - pops a path from the queue,
//    - checks if the path starts with a source asset, if it does, it appends
//      the path to results,
//    - finds all assets connected to the head of the current path and prepends
//      the path, calculating the current cost.
// Algorithm ends when there is no more paths to extend (len(queue) = 0) or
// `maxResults` has been reached.
//
// The actual calculation of the cost is happening in `pathNode.Cost()` method.
// There are a couple of important things to note:
// 1. We are given `DestinationAmount` and the destination asset is the tail of
// the list. So we need to start from the end of the path and continue to the
// front.
// 2. Because we are going from the tail to the head of the path, given the path
// A -> B -> C. we are interested in finding:
//    - amount of B needed to buy `DestinationAmount` of C,
//    - amount of A needed to by above amount of B.
// 3. Finally, the actual path payment will sell A, buy B, sell B, buy C. So we
// need to check the following orderbooks:
//    - sell C, buy B (the user will sell B, buy C),
//    - sell B, buy A (the user will sell A, buy B).
//
// The algorithm works as follows:
// 1. Because the head of the path is source account, we build a stack of assets
//    to reverse that order.
// 2. We start with the last asset (pop the stack), calculate it's cost (if not
//    cached) and continue towards the source asset (bottom of the stack).
// 3. We return the final cost.
package simplepath
