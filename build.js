#!/usr/bin/env node
const esbuild = require('esbuild');
const {sassPlugin} = require('esbuild-sass-plugin');
const {rmSync, readFileSync, writeFileSync} = require('fs');
const {basename} = require('path')
const glob = require('glob').sync;


//
// Args ...
//

const prod = process.argv.slice(2).includes('--prod');
const watch = process.argv.slice(2).includes('--watch')


//
// Cleanup first ...
//

rmSync('./pkg/webserver/static/assets', {
  recursive: true,
  force: true
});


//
// Go package ...
//

let monacoGoCode = `package monaco

// revive:disable:line-length-limit

// Languages ...
var Languages = []Language{
\t{
\t\tID:         "plain",
\t\tName:       "Plain text",
\t\tAliases:    []string{},
\t\tExtensions: []string{".txt"},
\t\tFilenames:  []string{},
\t},
`
glob('./node_modules/monaco-editor/esm/vs/basic-languages/*/*.contribution.js').forEach((file) => {
  const content = readFileSync(file).toString();
  const regexp = /registerLanguage(\(.*?^}\));/gms;
  const matches = content.matchAll(regexp);

  for (const match of matches) {
    let language = {}
    try {
      language = eval(match[1]);
    } catch (e) {
      console.log('--------\n'+match[1]+'\n--------');
      process.exit(0);
    }

    if (!language.aliases) {  // PLA?
      return;
    }

    let name = language.aliases[0];
    switch (name) {
      case 'abap':             /* fix case */   name = 'ABAP'; break;
      case 'eeStructuredText': /* do nothing */ break;
      case 'sol':              /* fix name */   name = 'Solidity'; break;
      case 'aes':              /* fix name */   name = 'Sophia'; break;
      case 'sparql':           /* fix case */   name = 'SPARQL'; break;
      default:                 /* capitalize */ name = name[0].toUpperCase() + name.slice(1); break;
    }
    language.extensions = language.extensions || [];
    language.filenames = language.filenames || [];

    monacoGoCode += `\t{
\t\tID:         "${language.id}",
\t\tName:       "${name}",
\t\tAliases:    []string{${JSON.stringify(language.aliases.slice(1)).replace(/,/g, ', ').slice(1, -1)}},
\t\tExtensions: []string{${JSON.stringify(language.extensions).replace(/,/g, ', ').slice(1, -1)}},
\t\tFilenames:  []string{${JSON.stringify(language.filenames).replace(/,/g, ', ').slice(1, -1)}},
\t},
`
  }
});
monacoGoCode += `}
`

writeFileSync('./pkg/monaco/languages.go', monacoGoCode);


//
// Monaco ...
//

esbuild.build({
  entryPoints: [
    './node_modules/monaco-editor/esm/vs/language/json/json.worker.js',
    './node_modules/monaco-editor/esm/vs/language/css/css.worker.js',
    './node_modules/monaco-editor/esm/vs/language/html/html.worker.js',
    './node_modules/monaco-editor/esm/vs/language/typescript/ts.worker.js',
    './node_modules/monaco-editor/esm/vs/editor/editor.worker.js',
  ],
  bundle: true,
  minify: prod,
  format: 'iife',
  outbase: './node_modules/monaco-editor/esm',
  outdir: './pkg/webserver/static/assets'
});


//
// Application
//

esbuild.build({
  entryPoints: [
    './web/app/app.js'
  ],
  bundle: true,
  minify: prod,
  format: 'iife',
  outbase: './web/app',
  outdir: './pkg/webserver/static/assets',
  loader: {
    '.ttf': 'file'
  },
  plugins: [sassPlugin()],
  watch: watch
});
