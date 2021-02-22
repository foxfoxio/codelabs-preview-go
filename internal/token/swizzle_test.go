package token

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSwizzle(t *testing.T) {
	testcases := []struct {
		in       string
		expected string
	}{
		{
			in:       "eyJhbGciOiJSUzI1NiIsImtpZCI6IjYxMDgzMDRiYWRmNDc1MWIyMWUwNDQwNTQyMDZhNDFkOGZmMWNiYTgiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiTmlwb24gQ2hpbmF0aGltYXRtb25na2hvbiIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS0vQU9oMTRHaHdBRWhpVU43UzFCV2RLbTFfLWFYWE9uTkwtU0lmM1VsRTJ0UmFydz1zOTYtYyIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS9mb3hmb3gtbGVhcm4iLCJhdWQiOiJmb3hmb3gtbGVhcm4iLCJhdXRoX3RpbWUiOjE2MTI4OTEzNzMsInVzZXJfaWQiOiJUSE1rWTFNZ1Y3YVBlM2dBYVpmMElENzB1UTgyIiwic3ViIjoiVEhNa1kxTWdWN2FQZTNnQWFaZjBJRDcwdVE4MiIsImlhdCI6MTYxMzI3MDQxMCwiZXhwIjoxNjEzMjc0MDEwLCJlbWFpbCI6Im5pcG9uLmNoaUBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJnb29nbGUuY29tIjpbIjExMzMwMzY0MDk0MDU3ODI4OTU2MSJdLCJlbWFpbCI6WyJuaXBvbi5jaGlAZ21haWwuY29tIl19LCJzaWduX2luX3Byb3ZpZGVyIjoiZ29vZ2xlLmNvbSJ9fQ.EQPSf7VuEOLuwxAMM_5K2X-gVfAM_KvX0-3X_UUiDg3coleEoh5aHx_ZmMmKE6zaXQSKqwxmovC-tagDtahXZTV4vYvt8u4qE3dx1vq6R4CnPO7Sf9Z2am2b_UrTmUja1x4eDaa4qzIRs14aZsgURmkPcI4K-6iaKxiNrNyPkuiAji6h-dHiX1-n9SSLO3yYADcWY0VvRKGYmvjI3kxp4mscOuObnSGiUGCQvGyw-r7B9mVv6Jss5zDtT2iX7iFWFcBCa5JpjpNTQOo80J28QfHC4p4iglMz-E7KHfZrgQSNzbrsqap5vaw7FfF8mgtxPLqH0N8k5V3pbJE-XuTckg",
			expected: "mxmYVRXWsYUSj9GNVZja50UMHJ3SRdXdm52NHV2R2gHY3RTcgJnczF2eNRFUoF3WllkS4QkL550aohWNvVjQJVmU3MTSxcDcOJ1UP92avtENiJ0QidkVHhGOXpWMVNXR5Zjc0l0N1dFb6EEOx5id6Z0dQRkRWhGSS9WYQRHUiRHb18WeqRDSrVnbYhkSTV3VvolVkNkQYpnMQtEVSpTbuATWol0Yuc2NotGQqRHbPpXTz1ka3xEYqVjLKVDSk9EbsNFVoJ3WgVDM0FlS5J3MiBWRkVzdyA2aU52UzRFYhNDbiFzW4cmU44UUtR0MTVjc1JzdlJjRwVDd5M3dYd3MXN1WXlGY1NEagVHLEVHcslndypEVQlFY7VjRK5GTulFY3lEY2cGcEZ2awJGNmVEaWRFYXRDLxc1dKBGTCV2Vm5yVzokNe5ETCdHe01kTGR3V2cmURBlRtIVZ6kEVhdXTutUb3NTW3hzMZpmbrhkeVhUWxlFNhpXQ0cldrNzV2NGWgtXSEtkOw0GS1hzMYZndYBWawMTWCtGSgtGNqF2dBlFY2lkeWdDSEFWcFhVYtlERLVWSUx0MUVlT1gURORDVFxUMqVETxg1eMhHT7xUeEtGSj92aIVHOzgldUhUYvhzMh9WS6R2NIpnYttGSjFXUvFWbRhFYqJ3bOpGVzIWaJhVW6tmbZpmdUllMJl0Y3gERZ12aulVcJlVWzgzRhFXRYFWbJR0Sq9yMhtGNEFWcFhVYvFkVgBXTutkd4gkYxRjbIdDSEFWcFhVYtlERLhHRFxUMitGT7R0aNlnbrhEenlVWqZHRMlHUFxENItHT5hVVMdDSENWar5GS0hkaMVDRXNGeiVUULF0aZJWRYB1bNVVWSV0MNh1YYNVeqJDYPdmRVpmbrhkaVRjYqZnaIpnZVRlMBtXTGtmRM52bXh1QjNDTtF0VYRDWyk1TFVlVzBjRSZVSq5kaQhFYnlUWZtXVvhEdMtXT7RUVOVDSVx0MEtmTqRFWhFXU0cFcRl1YplERLp2MuJWaVhUY1ZGNh52Z0EmbJpmTqBFWjlWSEtkaz4mYpVFShVnZ0EmbnRTYuhDVhdXTutUb3NTW3hzMZZ3MYl1c4g0YtlUWjtWVzI2d3omT7FUSjFzZuh0NIpnY7tmbIRHS6hVdYVlT7BzejpXRuRVMJVVU0VlMM52axQVd2x2U2hjRWpVRYt0ZFVVYNF1MVRUR7RFNzYVVxdGWRN0YJBWSRVFTwhjVQd3LUh1dvMTYrRDRjZXVINmd4MDW6V1MiJTVIF2b4MTYvRjeMB3dzs0dutnY4FVSjBXSq5kaU5mYyEFNYFXQvhEdIpWY3d2Mg9GNzEWdRlFW1tGSgFTRuFWcnNDUoNzMhh3auNlautGStBDWYZXS6R2LQdGaSBzVKtEaQhmQXZ2LLJUToh2Uah2TW5EbbZEUqd0QPd2WD5EeSN1T2J1QPZnVW5EeKZlTwQ2QPx2UWpFaTNkT5h2QOdnWppUNKJ0WvVHbKJnSo9EMKlnVStEaQhGZGN2ZLhnZ",
		},
	}

	for i, c := range testcases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			out := Swizzle(c.in)
			assert.Equal(t, c.expected, out)
		})
	}

}

func TestUnswizzle(t *testing.T) {
	testcases := []struct {
		in       string
		expected string
	}{
		{
			in:       "mxmYVRXWsYUSj9GNVZja50UMHJ3SRdXdm52NHV2R2gHY3RTcgJnczF2eNRFUoF3WllkS4QkL550aohWNvVjQJVmU3MTSxcDcOJ1UP92avtENiJ0QidkVHhGOXpWMVNXR5Zjc0l0N1dFb6EEOx5id6Z0dQRkRWhGSS9WYQRHUiRHb18WeqRDSrVnbYhkSTV3VvolVkNkQYpnMQtEVSpTbuATWol0Yuc2NotGQqRHbPpXTz1ka3xEYqVjLKVDSk9EbsNFVoJ3WgVDM0FlS5J3MiBWRkVzdyA2aU52UzRFYhNDbiFzW4cmU44UUtR0MTVjc1JzdlJjRwVDd5M3dYd3MXN1WXlGY1NEagVHLEVHcslndypEVQlFY7VjRK5GTulFY3lEY2cGcEZ2awJGNmVEaWRFYXRDLxc1dKBGTCV2Vm5yVzokNe5ETCdHe01kTGR3V2cmURBlRtIVZ6kEVhdXTutUb3NTW3hzMZpmbrhkeVhUWxlFNhpXQ0cldrNzV2NGWgtXSEtkOw0GS1hzMYZndYBWawMTWCtGSgtGNqF2dBlFY2lkeWdDSEFWcFhVYtlERLVWSUx0MUVlT1gURORDVFxUMqVETxg1eMhHT7xUeEtGSj92aIVHOzgldUhUYvhzMh9WS6R2NIpnYttGSjFXUvFWbRhFYqJ3bOpGVzIWaJhVW6tmbZpmdUllMJl0Y3gERZ12aulVcJlVWzgzRhFXRYFWbJR0Sq9yMhtGNEFWcFhVYvFkVgBXTutkd4gkYxRjbIdDSEFWcFhVYtlERLhHRFxUMitGT7R0aNlnbrhEenlVWqZHRMlHUFxENItHT5hVVMdDSENWar5GS0hkaMVDRXNGeiVUULF0aZJWRYB1bNVVWSV0MNh1YYNVeqJDYPdmRVpmbrhkaVRjYqZnaIpnZVRlMBtXTGtmRM52bXh1QjNDTtF0VYRDWyk1TFVlVzBjRSZVSq5kaQhFYnlUWZtXVvhEdMtXT7RUVOVDSVx0MEtmTqRFWhFXU0cFcRl1YplERLp2MuJWaVhUY1ZGNh52Z0EmbJpmTqBFWjlWSEtkaz4mYpVFShVnZ0EmbnRTYuhDVhdXTutUb3NTW3hzMZZ3MYl1c4g0YtlUWjtWVzI2d3omT7FUSjFzZuh0NIpnY7tmbIRHS6hVdYVlT7BzejpXRuRVMJVVU0VlMM52axQVd2x2U2hjRWpVRYt0ZFVVYNF1MVRUR7RFNzYVVxdGWRN0YJBWSRVFTwhjVQd3LUh1dvMTYrRDRjZXVINmd4MDW6V1MiJTVIF2b4MTYvRjeMB3dzs0dutnY4FVSjBXSq5kaU5mYyEFNYFXQvhEdIpWY3d2Mg9GNzEWdRlFW1tGSgFTRuFWcnNDUoNzMhh3auNlautGStBDWYZXS6R2LQdGaSBzVKtEaQhmQXZ2LLJUToh2Uah2TW5EbbZEUqd0QPd2WD5EeSN1T2J1QPZnVW5EeKZlTwQ2QPx2UWpFaTNkT5h2QOdnWppUNKJ0WvVHbKJnSo9EMKlnVStEaQhGZGN2ZLhnZ",
			expected: "eyJhbGciOiJSUzI1NiIsImtpZCI6IjYxMDgzMDRiYWRmNDc1MWIyMWUwNDQwNTQyMDZhNDFkOGZmMWNiYTgiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiTmlwb24gQ2hpbmF0aGltYXRtb25na2hvbiIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS0vQU9oMTRHaHdBRWhpVU43UzFCV2RLbTFfLWFYWE9uTkwtU0lmM1VsRTJ0UmFydz1zOTYtYyIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS9mb3hmb3gtbGVhcm4iLCJhdWQiOiJmb3hmb3gtbGVhcm4iLCJhdXRoX3RpbWUiOjE2MTI4OTEzNzMsInVzZXJfaWQiOiJUSE1rWTFNZ1Y3YVBlM2dBYVpmMElENzB1UTgyIiwic3ViIjoiVEhNa1kxTWdWN2FQZTNnQWFaZjBJRDcwdVE4MiIsImlhdCI6MTYxMzI3MDQxMCwiZXhwIjoxNjEzMjc0MDEwLCJlbWFpbCI6Im5pcG9uLmNoaUBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJnb29nbGUuY29tIjpbIjExMzMwMzY0MDk0MDU3ODI4OTU2MSJdLCJlbWFpbCI6WyJuaXBvbi5jaGlAZ21haWwuY29tIl19LCJzaWduX2luX3Byb3ZpZGVyIjoiZ29vZ2xlLmNvbSJ9fQ.EQPSf7VuEOLuwxAMM_5K2X-gVfAM_KvX0-3X_UUiDg3coleEoh5aHx_ZmMmKE6zaXQSKqwxmovC-tagDtahXZTV4vYvt8u4qE3dx1vq6R4CnPO7Sf9Z2am2b_UrTmUja1x4eDaa4qzIRs14aZsgURmkPcI4K-6iaKxiNrNyPkuiAji6h-dHiX1-n9SSLO3yYADcWY0VvRKGYmvjI3kxp4mscOuObnSGiUGCQvGyw-r7B9mVv6Jss5zDtT2iX7iFWFcBCa5JpjpNTQOo80J28QfHC4p4iglMz-E7KHfZrgQSNzbrsqap5vaw7FfF8mgtxPLqH0N8k5V3pbJE-XuTckg",
		},
	}

	for i, c := range testcases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			out := Unswizzle(c.in)
			assert.Equal(t, c.expected, out)
		})
	}

}
