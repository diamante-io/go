---
title: Errors
replacement: https://developers.diamcircle.org/api/errors/
---

In the event that an error occurs while processing a request to aurora, an
**error** response will be returned to the client.  This error response will
contain information detailing why the request couldn't complete successfully.

Like HAL for successful responses, aurora uses a standard to specify how we
communicate errors to the client.  Specifically, aurora uses the [Problem
Details for HTTP APIs](https://tools.ietf.org/html/draft-ietf-appsawg-http-problem-00) draft specification.  The specification is short, so we recommend
you read it.  In summary, when an error occurs on the server we respond with a
json document with the following attributes:

|   name   |  type  |                                                                        description                                                                        |
| -------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| type     | url    | The identifier for the error, expressed as a url.  Visiting the url in a web browser will redirect you to the additional documentation for the problem. |
| title    | string | A short title describing the error.                                                                                                                     |
| status   | number | An HTTP status code that maps to the error.  An error that is triggered due to client input will be in the 400-499 range of status code, for example.  |
| detail   | string | A longer description of the error meant the further explain the error to developers.                                                                   |
| instance | string | A token that uniquely identifies this request.  Allows server administrators to correlate a client report with server log files                           |


## Standard Errors

There are a set of errors that can occur in any request to aurora which we
call **standard errors**.  These errors are:

- [Server Error](../reference/errors/server-error.md)
- [Rate Limit Exceeded](../reference/errors/rate-limit-exceeded.md)
- [Forbidden](../reference/errors/forbidden.md)
