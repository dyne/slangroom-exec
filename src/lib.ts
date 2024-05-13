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
import { redis } from "@slangroom/redis";
import { shell } from "@slangroom/shell";
import { timestamp } from "@slangroom/timestamp";
import { wallet } from "@slangroom/wallet";
import { zencode } from "@slangroom/zencode";

import type { ZenParams } from "@slangroom/shared";
import type { StringKeyOf } from "type-fest";
type ZenParamKey = StringKeyOf<ZenParams>;

export const decode = (r: string) => Buffer.from(r, "base64").toString();
export const decode_and_trim = (r: string) => decode(r).trim();
export const decode_and_json = (r: string) => {
	try {
		return JSON.parse(decode(r));
	} catch (e) {
		console.error("JSON is malformed");
		process.exit(3);
	}
};

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
		redis,
		shell,
		timestamp,
		wallet,
		zencode,
	]);

	const [c, sl, d, k, e, cx] = input.split("\n");

	if (sl.trim().length === 0) {
		console.error("Slangroom contract is empty");
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
