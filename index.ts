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
import type { ZenParams } from "@slangroom/shared";
import { shell } from "@slangroom/shell";
import { timestamp } from "@slangroom/timestamp";
import { wallet } from "@slangroom/wallet";
import { zencode } from "@slangroom/zencode";

const the_input = await Bun.stdin.text();
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

const decode_and_trim = (r: string) =>
  Buffer.from(r, "base64").toString().trim();
const decode = (r: string) => Buffer.from(r, "base64").toString();
const decode_and_json = (r: string) => JSON.parse(decode(r));
type Names = "conf" | "data" | "keys" | "extra";

const decode_param = (source: string, name: Names, fn: Function) => {
  if (source && source !== "") {
    try {
      opts[name] = fn(source);
    } catch (e) {
      console.error(`${name} is malformed`);
      process.exit(2);
    }
  }
};

const [c, sl, d, k, e, cx] = the_input.split("\n");
const opts: ZenParams = { data: {}, keys: {} };
const contract = decode(sl);
decode_param(c, "conf", decode_and_trim);
decode_param(d, "data", decode_and_json);
decode_param(k, "keys", decode_and_json);
decode_param(e, "extra", decode_and_json);
const { result } = await s.execute(contract, opts);
await Bun.write(Bun.stdout, JSON.stringify(result));
