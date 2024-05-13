import { expect, test } from "bun:test";
import { encode } from "../src/lib";

test("the encode() utility should work ok", () => {
	const conf = "";
	const slang = `Rule unknown ignore
Given I fetch the local timestamp in seconds and output into 'timestamp'
Given I have a 'number' named 'timestamp'
Then print the 'timestamp'`;
	const have = encode(conf, slang, "", "", "");
	const want =
		"\nUnVsZSB1bmtub3duIGlnbm9yZQpHaXZlbiBJIGZldGNoIHRoZSBsb2NhbCB0aW1lc3RhbXAgaW4gc2Vjb25kcyBhbmQgb3V0cHV0IGludG8gJ3RpbWVzdGFtcCcKR2l2ZW4gSSBoYXZlIGEgJ251bWJlcicgbmFtZWQgJ3RpbWVzdGFtcCcKVGhlbiBwcmludCB0aGUgJ3RpbWVzdGFtcCc=\n\n\n\n";
	expect(have).toBe(want);
});
