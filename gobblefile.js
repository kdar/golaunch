var gobble = require( 'gobble' );
// var fs = require('fs');
// var path = require('path');
// var spawn = require('child_process').spawn;
//
// function compilePlugins(inputdir, options, callback) {
//   console.log(inputdir);
//   callback();
//   return;
//
//   var files = fs.readdirSync(options.baseDir);
//   for (var x = 0; x < files.length; x++) {
//     var dirPath = path.join(inputdir, files[x]);
//     var stats = fs.statSync(dirPath);
//     if (!stats.isDirectory()) {
//       continue;
//     }
//
//     var pluginPath = "./"+path.relative('./', dirPath);
//
//     var subfiles = fs.readdirSync(dirPath);
//     for (var y = 0; y < subfiles.length; y++) {
//       if (subfiles[y].endsWith(".go")) {
//         var child = spawn("go", ["build", "-o", path.join(pluginPath, path.basename(pluginPath)), pluginPath], {
//           stdio: [0, 1, 2]
//         });
//         child.on('close', function(code, signal) {
//           if (code != 0) {
//             callback(new Error("err"));
//           }
//         });
//         break;
//       }
//     }
//   }
//
//   callback();
// }

module.exports = gobble([
  gobble('src').transform('babel', {
    plugins: ["mjsx"]
  }),
  //gobble('plugins').include("*.go").observe(compilePlugins)
]);
