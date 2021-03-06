package main

import (
	"bytes"
	"encoding/base64"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	_ = os.Setenv("SITU_VAULT_PASSWORD", "test-pw")
	wd, _ := os.Getwd()

	os.Args = []string{"situ-vault-kustomize", "./testdata/example/secrets.yaml", wd + "/testdata/example/"}
	output, outputErr := captureOutput(main)

	assert.Equal(t, "", outputErr)
	assert.Contains(t, output, "kind: Secret\n")
	assert.Contains(t, output, "username: "+b64("test-data")+"\n")
	assert.Contains(t, output, "password: "+b64("test-data-longer")+"\n")

	assert.NotContains(t, output, "SITU_VAULT")
	assert.NotContains(t, output, ".yaml")
	assert.True(t, strings.HasSuffix(output, "---\n"))
	assert.True(t, strings.HasPrefix(output, expected))
}

func Test_getPasswordEnv(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	_ = os.Setenv("SITU_VAULT_PASSWORD_VAR", "test-pw")
	pwc := PasswordConfig{
		Env: "SITU_VAULT_PASSWORD_VAR",
	}
	password := getPassword(pwc)
	assert.Equal(t, "test-pw", password)
}

func Test_getPasswordFile(t *testing.T) {
	pwc := PasswordConfig{
		Env:  "UNDEFINED_VAR_NAME", // env var not actually present, should fallback to file
		File: "file://./test.key",
	}
	password := getPassword(pwc)
	assert.Equal(t, "file://./test.key", password)
}

// Helpers:

func b64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func captureOutput(function func()) (out string, err string) {
	var bufferOut bytes.Buffer
	var bufferErr bytes.Buffer
	logStdout.SetOutput(&bufferOut)
	logStderr.SetOutput(&bufferErr)
	function()
	return bufferOut.String(), bufferErr.String()
}

var expected = `apiVersion: v1
kind: Secret
metadata:
    name: secret-one
type: Opaque
stringData:
    greeting: Welcome
    host: example.org
data:
    password: dGVzdC1kYXRhLWxvbmdlcg==
    risky_flag: dHJ1ZQ==
    username: dGVzdC1kYXRh

---
apiVersion: v1
kind: Secret
metadata:
    name: secret-tls
type: kubernetes.io/tls
data:
    tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNwRENDQVl3Q0NRQ3Q0RHpTeUFydGdEQU5CZ2txaGtpRzl3MEJBUXNGQURBVU1SSXdFQVlEVlFRRERBbHMKYjJOaGJHaHZjM1F3SGhjTk1qRXdNVEEyTWpJeU1USTJXaGNOTWpJd01UQTJNakl5TVRJMldqQVVNUkl3RUFZRApWUVFEREFsc2IyTmhiR2h2YzNRd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUUM2CjFZblJDRDdaRlR5N21tNWFDb2VLMUdCYkhoTzdVblMweVk2bmtuZU0xWUo4cTM1VW5oRytzZ1dscSsyWWttVzcKV0J1bThNSWdmdm52RHZkMjA3WmhvSllxL1JrbWVxcTVyMXF3MUtPZFZ3bjJ4KytMWkl4ZjFJWEU4VEszaDUyNQpYRWZGWWt1Qi8rUHl2akVuZ0hvbTU0enVUYUZxaUg0bElQdzFaSWV0MGV5L29kU1JocWtMbVNBazV1L1Yvc3JNClp5ZjRIQjE2ODFETEl2Q1JYUWlsaTNJRW1WeVRHaFQxblhHSW5ZVTRBSlNrNjlVaExuS0xqN05VUm9NdDdsMlQKQ0tDUGxZbmkvblFOREppN1hwTi9wbDBRVGRpekM3aU9BK3cyZU9qZTRNaGkwMlllTnBXNjZwMTlFTnMrc0NXQQpHTW96TzhzVG4zcUx3YXA0SlBPWEFnTUJBQUV3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUZubW91VVFobS9uCkZ6Wjc5dGk0VmNJVU5oOFhsQlNtcHE5WlRlNGhlQUFmNHRIRmxRUmM3by9sSmdENTJENHkzLzRuOUVHcFJIa1oKVjJqVVM5SVR3VTR5d0locHQ0U0FpcVlHMC9NTkVUODQwNWJTbC9ucXFmcUwzSENKU3V0MUxqQkJCMWc4TzdIdwptYjJBQXRTV21vMUlGNUlJSzhicEtjZzFCM1MzbXhlOWliWVZIRzZBRTBqcUVtLzd4Mk5PRDI0VWRNTXFxZVl0CjhQU0FFdStnZW1PSDlwWUxsTDFtY2x2Y1dVRFUrN0lQSG5Ra0pNVk9tNktveUU2WDZxWHRtYW15ckN4bUR5N1oKcnFzME45enRodnNXNGZIaUFOSmZxdHNKSVlRTERiQXN1bXl1RkRrRHlDMEJnNDJrQU5lVXAyc1VsYkFlMFdkaAowRzExY21VTXp4ST0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    tls.key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2Z0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktnd2dnU2tBZ0VBQW9JQkFRQzYxWW5SQ0Q3WkZUeTcKbW01YUNvZUsxR0JiSGhPN1VuUzB5WTZua25lTTFZSjhxMzVVbmhHK3NnV2xxKzJZa21XN1dCdW04TUlnZnZudgpEdmQyMDdaaG9KWXEvUmttZXFxNXIxcXcxS09kVnduMngrK0xaSXhmMUlYRThUSzNoNTI1WEVmRllrdUIvK1B5CnZqRW5nSG9tNTR6dVRhRnFpSDRsSVB3MVpJZXQwZXkvb2RTUmhxa0xtU0FrNXUvVi9zck1aeWY0SEIxNjgxREwKSXZDUlhRaWxpM0lFbVZ5VEdoVDFuWEdJbllVNEFKU2s2OVVoTG5LTGo3TlVSb010N2wyVENLQ1BsWW5pL25RTgpESmk3WHBOL3BsMFFUZGl6QzdpT0ErdzJlT2plNE1oaTAyWWVOcFc2NnAxOUVOcytzQ1dBR01vek84c1RuM3FMCndhcDRKUE9YQWdNQkFBRUNnZ0VBZlpRZUNBUTB5aEMrTzVLM2JZbjZSTlF1MTgvRmoza0N2S2xsV3pqVlpqSDAKZlB0LzlEd3l6U3czSTM0R096RGJkQ3JxbXpEa0twZHVRc0thanFJS3lsLzN6M2xET0Z2bStOdm1aMGpsbUZIeQpmbzh3Y1U3cUUxZHpla1pzd25ORERsMzZWNitUOVJNY0VnTElZemExNUFScTg1bjJUdmJqWXUxaTJEaDBBZDRaCmtkMkxlcmZOZU9Kc2szR3k5THd1bzNpVzJHeE5Ydk5HNTdabkZNc2o4SzQybk52blE1dm91NkhhRlFNSGlyekQKL01DVzlOUzl0VUdZRVhCWWIzMzJIMVMvS0l3dzFKNnh2eTBBRkNodEo4SWdZVkgzbnBnWDd0K2dPeS9XSnQ0dwpXZk81Z1R3UnJKNUtmMXBsODhuVTlhYS8rMHlYQmdVOGZvcXFDcFppOFFLQmdRRDRiS3NzWUtZTjJsWUhGZzc5CjdOUy9NWFZDTm92RzdmWURBcnZnQWJpR1RLQzVEZDdlWjEvYXB1YjZ3MGcyUHFsK1hCMXMwN2JIazFldm1xQ1MKZHNZc3hPWG5xRDdxMWxjbkQya3RwZGtqM2NvZW82Yjl5NFNVcXVZTUl4SUZSK2daQlFLd3dQTkY3cWFvcFU4dwp5NGw3YkEvWHpBN2JEZS9QYWExeXhZVDEzUUtCZ1FEQWlCQXZiTUNlUWNkNkpaK1NBWEE2eGljVlRQRXBBUUlPCi9TRVdzK3JOUDN3S0VGMjlScGZuRUtvNTl3L3dNenU3eWhCa3pySUpSSzRUS3BhSFdqMDdGWTRNTXdWQVRldVkKcFlySUV2Nyt1VTd6czg2RSs1ZFVBMWF4eWFRV3VqN25zUU5zWm5LYzJBUkVSQkZHMzRiT2lBOFpBTlNRbXF1WQo3cnVXYmtvNkF3S0JnUUMyd290empHN2RoaUQvK1pSeDdzZmRHSytoVkt1a1gvQTY2c240MUl0Q0VpR3p3cWFSCmpBK1N0bkw3VEt0VmJPZ1kwLys1emsrTHA3UTh0azhuTVVZK0xXVE45cExEQllqOGJYUDlaeVBHSlNiTFA2NWMKekZydlhJTDlydGRWRnorREdKS1FJb05Xa1duK2JBOUVZSmoyT2R1MThLT0ZPRTJTazdaTTEwOG42UUtCZ1FDNgpBUG02QjVRVGtMTTUwNjFVN21UUnMyeEF6T1BUM0hCenNLTk4vcVhpZ3VuQUEwMjh0YjI5YzBFeDNQbWQ4ckZMCjNJeDRCNlREQllJemJCcWZTMVFLaCttQzZhdXlFMVdBVkxZK1V2UGRmWVBFTjd0V2lJWUxtV29oT3hCM0VKb0QKVnVWYXphTCsya2RNK0lIRWVlRVFHU3lVMkZPRUhKbVpsMUxObzJHOHB3S0JnRHhhRWg3ZEEyaDJlckMrZGhGawo5a29zZml3TmVJSzhDRzk5STlXQVhEU2JyNjV4TjA5aEx1QzBMNFg1WmdGeVJ2bTRGUE4vdHBSZzJvLzJHZC95Cm9BTUhmbGFYYzVmc0lLOVZSSDhRT2tRSlFOcjBoeVRFY1g2Rmc2TytkMm0vV243MlAyS3VOOEh1NXBZQ1RqbzIKckFHRXNYdURrNWdSZ04wZ2V5eTZLRkZyCi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K

---
`
