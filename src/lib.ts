// SPDX-FileCopyrightText: 2024-2025 Dyne.org foundation
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Slangroom } from "@slangroom/core";
import { ethereum } from "@slangroom/ethereum";
import { fs } from "@slangroom/fs";
import { git } from "@slangroom/git";
import { helpers } from "@slangroom/helpers";
import { http } from "@slangroom/http";
import { JSONSchema } from "@slangroom/json-schema";
import { oauth } from "@slangroom/oauth";
import { pocketbase } from "@slangroom/pocketbase";
import { qrcode } from "@slangroom/qrcode";
import { rdf } from "@slangroom/rdf";
import { redis } from "@slangroom/redis";
import { shell } from "@slangroom/shell";
import { timestamp } from "@slangroom/timestamp";
import { wallet } from "@slangroom/wallet";
import { zencode } from "@slangroom/zencode";
import { execute } from "@dyne/slangroom-chain";

import type { ZenParams } from "@slangroom/shared";
type ZenParamKey = Extract<keyof ZenParams, string>;

export const decode = (r: string) => Buffer.from(r, "base64").toString();
export const decode_and_trim = (r: string) => decode(r).trim();
export const decode_and_json = (r: string) => JSON.parse(decode(r));

export const slangroom_exec = async (input: string) => {
	const s = new Slangroom([
		ethereum,
		fs,
		git,
		helpers,
		http,
		JSONSchema,
		oauth,
		pocketbase,
		qrcode,
		rdf,
		redis,
		shell,
		timestamp,
		wallet,
		zencode,
	]);

	const [c, sl, d, k, e, cx] = input.split("\n");

	if (!sl || sl.trim().length === 0) {
		console.error("Malformed input: Slangroom contract is empty");
		process.exit(1);
	}

	const opts: ZenParams = { data: {}, keys: {} };
	const contract = decode(sl);

	const decode_slangroom_param = (
		source: string,
		key: ZenParamKey,
		fn: Function
	) => {
		if (source && source !== "") {
			try {
				opts[key] = fn(source);
			} catch (e) {
				console.error(`${key.toUpperCase()} is malformed`);
				process.exit(2);
			}
		}
	};

	decode_slangroom_param(c, "conf", decode_and_trim);
	decode_slangroom_param(d, "data", decode_and_json);
	decode_slangroom_param(k, "keys", decode_and_json);
	decode_slangroom_param(e, "extra", decode_and_json);
	const { result } = await s.execute(contract, opts);
	return JSON.stringify(result);
};

export const encode = (
	conf?: string,
	contract?: string,
	data?: string,
	keys?: string,
	extra?: string,
	ctx?: string
) => {
	const b64 = (source: string) => Buffer.from(source).toString("base64");

	return [
		b64(conf ?? ""),
		b64(contract ?? ""),
		b64(data ?? ""),
		b64(keys ?? ""),
		b64(extra ?? ""),
		b64(ctx ?? ""),
	].join("\n");
};

export const slangroom_chain_exec = async (input: string) => {
	// matain same input, but only contract and data are used to run a chain
	const [_, sl, d] = input.split("\n");
	const chain = decode(sl);
	const data = decode(d);
	const result = await execute(chain, data);
	return result;
}
