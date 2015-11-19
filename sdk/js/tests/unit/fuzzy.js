define(function (require) {
  var registerSuite = require('intern!object');
  var assert = require('intern/chai!assert');
  var fuzzy = require('intern/dojo/node!../../fuzzy');

  registerSuite(function() {
    // variables go here

    return {
      name: 'fuzzy',

      'Test match': function () {
        var tableTests = [
        	["S", "Shutdown", 76],
        	["h", "Shutdown", 60],
        	["u", "Shutdown", 50],
        	["t", "Shutdown", 43],
        	["d", "Shutdown", 38],
        	["o", "Shutdown", 35],
        	["w", "Shutdown", 32],
        	["n", "Shutdown", 30],
        	["chrome", "chrome", 107],
        	["chrom", "chrome", 105],
        	["chrom", "chrom", 105],
        	["nix", "i love me some unix", 19],
        	["nix", "unix i love", 76],
        	["abc", "def", 0],
        	["abc", "zxyabdef", 0],
        	["ad", "zxyabdef", 47],
          ["", "hey", 0],
          ["hey", "", 0],
        ];

        for (var x = 0; x < tableTests.length; x++) {
          assert.strictEqual(fuzzy.match(tableTests[x][0], tableTests[x][1]).score, tableTests[x][2], "test #" + x);
        }
      }
    };
  });
});
