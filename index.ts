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

for await (const chunk of Bun.stdin.stream()) {
  // chunk is Uint8Array
  // this converts it to text (assumes ASCII encoding)
  const the_input = Buffer.from(chunk).toString();
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

  const { result } = await s.execute(the_input, {});
  console.log(JSON.stringify(result));
}
