package main

import "testing"

func TestConcat(t *testing.T) {
	cases := []struct {
		elm  string
		elms []string
		want []string
	}{
		{
			elm:  "",
			elms: []string(nil),
			want: []string{""},
		},
		{
			elm:  "owl",
			elms: []string{"fox"},
			want: []string{"owl", "fox"},
		},
		{
			elm:  "owl",
			elms: []string{"fox", "raccoon"},
			want: []string{"owl", "fox", "raccoon"},
		},
	}

	for _, c := range cases {
		ret := concat(c.elm, c.elms)

		if len(ret) != len(c.want) {
			t.Fatalf("size unmatch r: %q want: %q\n", ret, c.want)
		}

		for i, r := range ret {
			if r != c.want[i] {
				t.Fatalf("ret: %q want %q\n", ret, c.want)
			}
		}
	}
}
