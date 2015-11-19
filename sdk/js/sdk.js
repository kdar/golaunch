var path = require('path');
var fs = require('fs');

module.exports = {
  imageFileToEmbed: function(path_) {
    return new Promise(function(resolve, reject) {
      fs.readFile(path_, function(err, data) {
        if (err) {
          reject(Error(err));
        } else {
          var base64data = new Buffer(data).toString('base64');
          var ext = path.extname(path_);
          if (ext) {
            ext = ext.substr(1);
          } else {
            ext = "png";
          }
          resolve("data:image/" + ext + ";base64," + base64data);
        }
      });
    });
  },
  imageFileToEmbedSync: function(path_) {
    var data = fs.readFileSync(path_);
    var base64data = new Buffer(data).toString('base64');
    var ext = path.extname(path_);
    if (ext) {
      ext = ext.substr(1);
    } else {
      ext = "png";
    }
    return "data:image/" + ext + ";base64," + base64data;
  }
};
