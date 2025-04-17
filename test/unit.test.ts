import { describe, expect, test } from "bun:test";
import {
	decode,
	decode_and_json,
	decode_and_trim,
	encode,
	slangroom_exec,
	slangroom_chain_exec,
} from "../src/lib";

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

test.each([
	// Empty string
	["", ""],
	// String with special characters
	["IUAjJCVeJiooKQ==", "!@#$%^&*()"],
	// String with Unicode characters
	["8J+agPCfjJ/wn4yI", "🚀🌟🌈"],
	// String with numbers
	["MTIzNDU2Nzg5MA==", "1234567890"],
	// String with spaces
	["   ", ""],
	// String with newline characters
	["bGluZTEKbGluZTIKbGluZTM=", "line1\nline2\nline3"],
	// String with tabs
	["CQkJCQ==", "\t\t\t\t"],
	// String with special characters and Unicode characters
	["JOKCrMKjwqXigrnwn42O8J+Ni/CfjZM=", "$€£¥₹🍎🍋🍓"],
	// String with a mix of uppercase and lowercase characters
	["QWJDZEVmR2hJaktsTW5PcFFyU3RVdld4WXo", "AbCdEfGhIjKlMnOpQrStUvWxYz"],
	// String with control characters
	["AAECAwQF", "\u0000\u0001\u0002\u0003\u0004\u0005"],
	// String with control characters
	["AAECAwQF", "\x00\x01\x02\x03\x04\x05"],
	// String with extended ASCII characters
	["w4DDgcOCw4PDhMOFw4bDh8OIw4nDisOLw4zDjcOOw48=", "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏ"],
])("the decode() should work correctly", (source, want) => {
	// runs once for each test case provided
	const have = decode(source);
	expect(have).toBe(want);
});

[
	// Empty string
	["", ""],
	// String with special characters
	["IUAjJCVeJiooKQ==", "!@#$%^&*()"],
	// String with Unicode characters
	["8J+agPCfjJ/wn4yI", "🚀🌟🌈"],
	// String with numbers
	["MTIzNDU2Nzg5MA==", "1234567890"],
	// String with spaces
	["   ", ""],
	// String with newline characters
	["bGluZTEKbGluZTIKbGluZTM=", "line1\nline2\nline3"],
	// String with tabs
	["CQkJCQ==", "\t\t\t\t"],
	// String with special characters and Unicode characters
	["JOKCrMKjwqXigrnwn42O8J+Ni/CfjZM=", "$€£¥₹🍎🍋🍓"],
	// String with a mix of uppercase and lowercase characters
	["QWJDZEVmR2hJaktsTW5PcFFyU3RVdld4WXo", "AbCdEfGhIjKlMnOpQrStUvWxYz"],
	// String with control characters
	["AAECAwQF", "\u0000\u0001\u0002\u0003\u0004\u0005"],
	// String with control characters
	["AAECAwQF", "\x00\x01\x02\x03\x04\x05"],
	// String with extended ASCII characters
	["w4DDgcOCw4PDhMOFw4bDh8OIw4nDisOLw4zDjcOOw48=", "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏ"],
];

test.each([
	// A simple newline
	["\n", ""],
	// Empty string
	["", ""],
	// String with special characters
	["   IUAjJCVeJiooKQ==   ", "!@#$%^&*()"],
	// String with Unicode characters
	["8J+agPCfjJ/wn4yI", "🚀🌟🌈"],
	// String with numbers
	["MTIzNDU2Nzg5MA==", "1234567890"],
	// String with spaces
	["   ", ""],
	// String with newline characters
	["bGluZTEKbGluZTIKbGluZTM=", "line1\nline2\nline3"],
	// String with encoded tabs
	["CQkJCQ==", ""],
	// String with special characters and Unicode characters
	["JOKCrMKjwqXigrnwn42O8J+Ni/CfjZM=", "$€£¥₹🍎🍋🍓"],
	// String with a mix of uppercase and lowercase characters
	["QWJDZEVmR2hJaktsTW5PcFFyU3RVdld4WXo", "AbCdEfGhIjKlMnOpQrStUvWxYz"],
	// String with control characters
	["AAECAwQF", "\u0000\u0001\u0002\u0003\u0004\u0005"],
	// String with control characters
	["AAECAwQF", "\x00\x01\x02\x03\x04\x05"],
	// String with extended ASCII characters
	["w4DDgcOCw4PDhMOFw4bDh8OIw4nDisOLw4zDjcOOw48=", "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏ"],
])("the decode_and_trim() should work correctly", (source, want) => {
	// runs once for each test case provided
	const have = decode_and_trim(source);
	expect(have).toBe(want);
});

// generate test for decode_and_json

describe("decode_and_json", () => {
	test.each([
		// String with special characters
		// ["eyJzb3VyY2UiOiAiJCTlwrzCvMK9wr7Cv8K/In0=", { source: "$€£¥₹🍎🍋🍓" }],
		// String with a mix of uppercase and lowercase characters
		[
			"eyJzb3VyY2UiOiAiQWJDZEVGaEppS2xNbk9wUXJTdFV2V3hZeiJ9",
			{ source: "AbCdEFhJiKlMnOpQrStUvWxYz" },
		],
		// String with control characters
		[
			"eyJzb3VyY2UiOiAiXHUwMDAwXHUwMDAxXHUwMDAyXHUwMDAzXHUwMDA0XHUwMDA1In0=",
			{ source: "\u0000\u0001\u0002\u0003\u0004\u0005" },
		],
		// String with control characters
		// [
		// 	"eyJzb3VyY2UiOiAiXDAwXDAxXDAyXDAzXDA0XDA1In0=",
		// 	{ source: "\x00\x01\x02\x03\x04\x05" },
		// ],
		// String with extended ASCII characters
		[
			"eyJzb3VyY2UiOiAiw4DDgcOCw4PDhMOFw4bDh8OIw4nDisOLw4zDjcOOw48ifQ==",
			{ source: "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏ" },
		],
	])(
		"should correctly decode and parse JSON with strings",
		(encoded, expected) => {
			const result = decode_and_json(encoded);
			expect(result).toEqual(expected);
		}
	);

	// test.each([
	// 	// Property name with special characters
	// 	[
	// 		"eyJAJCXCvMK9wr7Cv8K/In0iOiAiJCTlwrzCvMK9wr7Cv8K/In0=",
	// 		{ "@$€£¥₹🍎🍋🍓": "$€£¥₹🍎🍋🍓" },
	// 	],
	// 	// Property name with a mix of uppercase and lowercase characters
	// 	[
	// 		"eyJBYkNkRWZHaEpqS2xNbk9wUXJTdFV2V3hZeiI6ICJBQkNERUZHaEpqS2xNbk9wUXJTdFV2V3hZeiI=",
	// 		{ AbCdEfGhIjKlMnOpQrStUvWxYz: "AbCdEfGhIjKlMnOpQrStUvWxYz" },
	// 	],
	// 	// Property name with control characters
	// 	[
	// 		"eyJcdTAwMDBcdTAwMDFcdTAwMDJcdTAwMDNcdTAwMDRcdTAwMDUiOiAiXHUwMDAwXHUwMDAxXHUwMDAyXHUwMDAzXHUwMDA0XHUwMDA1In0=",
	// 		{
	// 			"\u0000\u0001\u0002\u0003\u0004\u0005":
	// 				"\u0000\u0001\u0002\u0003\u0004\u0005",
	// 		},
	// 	],
	// 	// Property name with control characters
	// 	[
	// 		"eyJcMDBcMDFcMDJcMDNcMDRcMDUiOiAiXDAwXDAxXDAyXDAzXDA0XDA1In0=",
	// 		{ "\x00\x01\x02\x03\x04\x05": "\x00\x01\x02\x03\x04\x05" },
	// 	],
	// 	// Property name with extended ASCII characters
	// 	[
	// 		"eyJ3w4DDgcOCw4PDhMOFw4bDh8OIw4nDisOLw4zDjcOOw48iOiAiw4DDgcOCw4PDhMOFw4bDh8OIw4nDisOLw4zDjcOOw48ifQ==",
	// 		{ ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏ: "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏ" },
	// 	],
	// ])(
	// 	"should correctly decode and parse JSON with charpame",
	// 	(encoded, expected) => {
	// 		const result = decode_and_json(encoded);
	// 		expect(result).toEqual(expected);
	// 	}
	// );

	test.each([
		// ...
		// Null value
		["eyJzb3VyY2UiOiBudWxsfQ==", { source: null }],
		// Date value
		[
			"eyJzb3VyY2UiOiAiMjAyMi0wMS0wMVQwMDowMDowMC4wMDBaIn0=",
			{ source: "2022-01-01T00:00:00.000Z" },
		],
		// Time value
		["eyJzb3VyY2UiOiAiMTI6MzQ6NTYifQ==", { source: "12:34:56" }],
		// Array value
		["eyJzb3VyY2UiOiBbMSwgMiwgMywgNCwgNV19", { source: [1, 2, 3, 4, 5] }],
	])(
		"should correctly decode and parse JSON with different values",
		(encoded, expected) => {
			const result = decode_and_json(encoded);
			expect(result).toEqual(expected);
		}
	);
});
describe("slangroom_exec", () => {
	test.each([
		[encode("", "Given nothing\nThen print data", "", "", ""), "[]"],
		[
			encode("", "Given nothing\nThen print the string 'ciao'", "", "", ""),
			'{"output":["ciao"]}',
		],
		[
			encode(
				"",
				"Given I have a 'string' named 'mimmo'\nThen print the 'mimmo'",
				'{"mimmo": "hello"}',
				"",
				""
			),
			'{"mimmo":"hello"}',
		],
		// [encode("", "Given nothing", "", "", ""), '{"mimmo":"hello"}'],
	])(`the slangroom_exec() should work correctly`, async (input, want) => {
		const have = await slangroom_exec(input);
		expect(have).toBe(want);
	});

	test.each([[""], ["Given nothing"], ["\\n"]])(
		"should throw an error if the Slangroom contract is empty",
		async (input) => {
			const proc = Bun.spawn(
				["bun", "-e", `require("./src/lib").slangroom_exec('${input}')`],
				{ stderr: "pipe" }
			);

			const error = await new Response(proc.stderr).text();
			await proc.exited;
			expect(error).toContain("Slangroom contract is empty");
			expect(proc.exitCode).toBe(1);
		}
	);

	test("should throw an error if the data is malformed", async () => {
		const i = encode("", "Given", "aspoks", "", "", "").replaceAll("\n", "\\n");
		const proc = Bun.spawn(
			["bun", "-e", `require("./src/lib").slangroom_exec('${i}')`],
			{ stderr: "pipe" }
		);
		await proc.exited;
		expect(proc.exitCode).toBe(2);

		const error = await new Response(proc.stderr).text();
		expect(error).toContain("DATA is malformed");
	});

	test("should throw an error if the keys are malformed", async () => {
		const i = encode("", "Given", "", "aspoks", "", "").replaceAll("\n", "\\n");
		const proc = Bun.spawn(
			["bun", "-e", `require("./src/lib").slangroom_exec('${i}')`],
			{ stderr: "pipe" }
		);
		await proc.exited;
		expect(proc.exitCode).toBe(2);

		const error = await new Response(proc.stderr).text();
		expect(error).toContain("KEYS is malformed");
	});

	test("should throw an error if the extra is malformed", async () => {
		const i = encode("", "Given", "", "", "aspoks", "").replaceAll("\n", "\\n");
		const proc = Bun.spawn(
			["bun", "-e", `require("./src/lib").slangroom_exec('${i}')`],
			{ stderr: "pipe" }
		);
		await proc.exited;
		expect(proc.exitCode).toBe(2);

		const error = await new Response(proc.stderr).text();
		expect(error).toContain("EXTRA is malformed");
	});
});

describe("slangroom_chain_exec", () => {
	const chain = `
steps:
  - id: hello
    zencode: |
      Given I have a 'string' named 'hello'
      Then print the 'hello'
      Then print the string 'Hello, world!'`;
	test.each([
		[
			encode("", chain, '{"hello":"hello from data!"}', "", ""),
			'{"hello":"hello from data!","output":["Hello,_world!"]}',
		],
	])(`the slangroom_chain_exec() should work correctly`, async (input, want) => {
		const have = await slangroom_chain_exec(input);
		expect(have).toBe(want);
	});
});
