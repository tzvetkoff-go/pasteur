const fs = require('fs');
const path = require('path');
const sass = require('node-sass');

module.exports = function(grunt) {
  var sassCommon;
  var uglifyCommon;

  grunt.initConfig({
    pkg: grunt.file.readJSON('package.json'),

    //
    // Concat all stylesheets
    //

    concat: {
      sass: {
        src: [
          // CodeMirror + theme + addon/mode styles
          './node_modules/codemirror/lib/codemirror.css',
          './node_modules/codemirror/theme/material-darker.css',
          './node_modules/codemirror/addon/search/matchesonscrollbar.css',
          './node_modules/codemirror/mode/*/*.css',
          // CodeMirror whitespaces plugin
          './node_modules/codemirror-whitespaces/dist/codemirror-whitespaces.css',
          // Select2
          './node_modules/select2/dist/css/select2.css',
          // Select2 dark theme
          './websrc/select2-dark-theme/select2-dark-theme.scss',
          // App
          './websrc/app/app.scss',
        ],
        dest: './tmp/app.bundle.css',
        nonull: true
      },
    },

    //
    // Build SASS/SCSS
    //

    sass: {
      common: (sassCommon = {
        src: [
          './tmp/app.bundle.css'
        ],
        dest: './pkg/webserver/static/assets/app.bundle.css',
        nonull: true,
      }),

      dist: Object.assign({}, sassCommon, {
        options: {
          implementation: sass,
          outputStyle: 'compressed'
        }
      }),

      dev: Object.assign({}, sassCommon, {
        options: {
          implementation: sass,
          outputStyle: 'expanded'
        }
      })
    },

    //
    // Build JS
    //

    uglify: {
      common: (uglifyCommon = {
        src: [
          // jQuery
          './node_modules/jquery/dist/jquery.slim.js',
          // CodeMirror + addons + modes
          './node_modules/codemirror/lib/codemirror.js',
          './node_modules/codemirror/addon/edit/matchbrackets.js',
          './node_modules/codemirror/addon/edit/closebrackets.js',
          './node_modules/codemirror/addon/edit/matchtags.js',
          './node_modules/codemirror/addon/edit/closetag.js',
          './node_modules/codemirror/addon/mode/simple.js',
          './node_modules/codemirror/addon/mode/multiplex.js',
          './node_modules/codemirror/addon/mode/overlay.js',
          './node_modules/codemirror/addon/search/search.js',
          './node_modules/codemirror/addon/search/searchcursor.js',
          './node_modules/codemirror/addon/search/jump-to-line.js',
          './node_modules/codemirror/addon/search/matchesonscrollbar.js',
          './node_modules/codemirror/addon/selection/active-line.js',
          './node_modules/codemirror/mode/meta.js',
          './node_modules/codemirror/mode/*/*.js',
          // CodeMirror whitespaces plugin
          './node_modules/codemirror-whitespaces/dist/codemirror-whitespaces.js',
          // Select2
          './node_modules/select2/dist/js/select2.full.js',
          // App
          './websrc/app/app.js'
        ],
        dest: './pkg/webserver/static/assets/app.bundle.js',
        nonull: true,
      }),

      dist: Object.assign({}, uglifyCommon, {
        options: {
          banner: '/*! Pasteur */ ',
          output: {
            beautify: false
          },
          mangle: true,
          compress: true
        }
      }),

      dev: Object.assign({}, uglifyCommon, {
        options: {
          banner: '/* Pasteur-DEV */',
          output: {
            beautify: true
          },
          mangle: false,
          compress: false
        }
      })
    },

    //
    // Watch
    //
    watch: {
      sass: {
        files: ['websrc/**/*.scss'],
        tasks: ['concat:sass', 'sass:dev']
      },
      uglify: {
        files: ['websrc/**/*.js'],
        tasks: ['uglify:dev']
      }
    }
  });

  //
  // Load tasks from libraries
  //

  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-sass');

  //
  // Generate pkg/codemirror/modes.go
  //

  grunt.registerTask('go:generate', 'Generate pkg/codemirror/modes.go', function() {
    var src = fs.readFileSync(path.resolve(__dirname, 'node_modules', 'codemirror', 'mode', 'meta.js')).toString();
    var frag = src.match(/CodeMirror\.modeInfo = (\[.*?\];)/ms)[1];
    var modes = eval(frag);
    var plainIndex = -1;

    for (var i = 0, len = modes.length; i < len; ++i) {
      var info = modes[i];
      if (info.mime == 'text/plain' || (info.mimes && info.mimes.indexOf('text/plain') !== -1)) {
        plainIndex = i;
        break;
      }
    }

    if (plainIndex !== -1) {
      modes.unshift(modes.splice(plainIndex, 1)[0]);
    }

    modes.unshift({ name: 'Auto', mode: 'null', mimes: [''] });

    var codeMirrorModesGoCode = `package codemirror

// revive:disable:line-length-limit

// Modes ...
var Modes = []Mode{
`
    for (var i = 0, len = modes.length; i < len; ++i) {
      var info = modes[i];
      var name = info.name;
      var mode = info.mode;
      var mimeTypes = info.mimes ? info.mimes : info.mime ? [info.mime] : [];
      var extensions = info.ext ? info.ext : [];
      var aliases = info.alias ? info.alias : [];
      codeMirrorModesGoCode += `\t{
\t\tName:       "${name}",
\t\tMode:       "${mode}",
\t\tMimeTypes:  []string{${JSON.stringify(mimeTypes).slice(1, -1).replace(/","/g, '", "')}},
\t\tExtensions: []string{${JSON.stringify(extensions).slice(1, -1).replace(/","/g, '", "')}},
\t\tAliases:    []string{${JSON.stringify(aliases).slice(1, -1).replace(/","/g, '", "')}},
\t},
`
    }
    codeMirrorModesGoCode += `}
`
    fs.writeFileSync(path.resolve(__dirname, 'pkg', 'codemirror', 'modes.go'), codeMirrorModesGoCode);
  });

  //
  // Distribution build
  //

  grunt.registerTask('default', [
    'go:generate',
    'concat:sass',
    'sass:dist',
    'uglify:dist'
  ]);

  //
  // Development build
  //

  grunt.registerTask('dev', [
    'go:generate',
    'concat:sass',
    'sass:dev',
    'uglify:dev'
  ]);

  //
  // Development build + watch
  //

  grunt.registerTask('watch', [
    'dev',
    'watch'
  ])
};
