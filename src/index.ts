// SPDX-FileCopyrightText: 2024-2025 Dyne.org foundation
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { slangroom_chain_exec, slangroom_exec } from "./lib";
import { introspect } from 'zenroom';
const packageJson = require("../package.json");

const argv = Bun.argv;

if (argv.includes("-v") || argv.includes("--version")) {
	console.log(`${packageJson.name} ${packageJson.version}`);
  console.error(`
      _
     | |
  ___| | _____ ____   ____  ____ ___   ___  ____ _____ _____ _   _ _____  ____
 /___) |(____ |  _ \\ / _  |/ ___) _ \\ / _ \\|    (_____) ___ ( \\ / ) ___ |/ ___)
|___ | |/ ___ | | | ( (_| | |  | |_| | |_| | | | |    | ____|) X (| ____( (___
(___/ \\_)_____|_| |_|\\___ |_|   \\___/ \\___/|_|_|_|    |_____|_/ \\_)_____)\\____)
                    (_____|
`)
	console.error(`Built for ${process.platform}-${process.arch}
Copyright (C) 2024-2025 ${packageJson.author}
License AGPL-3.0-or-later: GNU AGPL version 3 <https://www.gnu.org/licenses/agpl-3.0.html>
This is free software: you are free to change and redistribute it
There is NO WARRANTY, to the extent permitted by law.`);
	process.exit(0);
}

if (argv.includes("-h") || argv.includes("--help")) {
	console.log(`Usage: cat <encoded_input> | ${packageJson.name} [options]
Options:
  -v, --version    Show version
  -h, --help       Show this help message
  -c, --chain      Execute the slangroom chain
  -i, --introspect Output the zenroom introspection of the input`);
	process.exit(0);
}

const the_input = await Bun.stdin.text();

if (argv.includes("-i")) {

	const the_output = await introspect(the_input);

	await Bun.write(Bun.stdout, JSON.stringify(the_output) + "\n");
} else if (argv.includes("-c")) {
	const the_output = await slangroom_chain_exec(the_input);
	await Bun.write(Bun.stdout, the_output);
} else {
	const the_output = await slangroom_exec(the_input);
	await Bun.write(Bun.stdout, the_output);
}
