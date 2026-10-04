package password

import "testing"

var fast = Params{Memory: 8 * 1024, Iterations: 1, Parallelism: 1, SaltLength: 16, KeyLength: 32}

func TestHashVerify(t *testing.T) {
	h, err := HashWithParams("correct horse battery", fast)
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := Verify("correct horse battery", h); err != nil || !ok {
		t.Fatalf("verify: %v %v", ok, err)
	}
	if ok, _ := Verify("wrong", h); ok {
		t.Fatal("wrong password accepted")
	}
	h2, _ := HashWithParams("correct horse battery", fast)
	if h == h2 {
		t.Fatal("salts must differ")
	}
	if _, err := HashWithParams("short", fast); err == nil {
		t.Fatal("short password must be rejected")
	}
	if _, err := Verify("x", "$argon2i$v=19$m=1,t=1,p=1$abc$def"); err == nil {
		t.Fatal("non-argon2id hash must be rejected")
	}
	if _, err := Verify("x", "garbage"); err == nil {
		t.Fatal("garbage must be rejected")
	}
}
