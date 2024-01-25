# Recovery Signer: Firebase Setup

This service uses Firebase to authenticate a user with an email address or
phone number. To configure a new Firebase project for use with recoverysigner
follow the steps below. These steps assume a default Firebase project setup.

## Enable Phone Number Authentication

1. Login to Firebase Console.
2. Click `Authentication` under `Develop` in the left sidebar.
3. Click `Sign-in method`.
4. Click `Phone`.
5. Toggle the feature to `Enable`.
6. Click `Save`.

## Enable Email Address Authentication

1. Login to Firebase Console.
2. Click `Authentication` under `Develop` in the left sidebar.
3. Click `Sign-in method`.
4. Click `Email/Password`.
5. Toggle the feature to `Enable`.
6. Toggle the `Email link (passwordless sign-in)`.
7. Click `Save`.
