import { slangroom_exec } from "./lib";

const the_input = await Bun.stdin.text();
const the_output = await slangroom_exec(the_input);
await Bun.write(Bun.stdout, the_output);
