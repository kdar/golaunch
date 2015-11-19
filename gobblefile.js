var gobble = require( 'gobble' );
module.exports = gobble([
  gobble('src').transform('babel', {
    plugins: ["mjsx"]
  })
]);
