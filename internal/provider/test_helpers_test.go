package provider

// Shared test fixtures. These constants exist purely to keep credential-shaped
// strings out of literal test configs so static-analysis credential scanners
// don't flag the repo. Nothing here is or has ever been a real credential.

const (
	// testFixturePassword is used wherever a test config requires a non-empty
	// admin password. The provider stores it in state as opaque text — it is
	// never validated against a real auth backend.
	testFixturePassword = "fixture-not-real-do-not-use"

	// testFixtureSecretValue is used for Key Vault Secret values in tests.
	testFixtureSecretValue = "fixture-secret-not-real"

	// testFixtureStorageKey is used for Function App storage_account_access_key
	// in tests. It is intentionally not a valid base64 SAS key.
	testFixtureStorageKey = "fixture-storage-key-not-real"
)
