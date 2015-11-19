function match(query, str) {
  if (!query || !str || query.length == 0 || str.length == 0) {
    return {
      success: false,
      score: 0
    };
  }

  var strLen = str.length;
  var compareString = str.toLowerCase();
  var pattern = query.toLowerCase();

  var sb = "";
  var patternIdx = 0;
  var firstMatchIndex = -1;
  var lastMatchIndex = 0;
  var ch;

  for (var idx = 0; idx < strLen; idx++) {
    ch = str[idx];
    if (compareString[idx] == pattern[patternIdx]) {
      if (firstMatchIndex < 0) {
        firstMatchIndex = idx;
      }
      lastMatchIndex = idx + 1;

      sb += /*prefix + */ ch /*+ suffix*/;
      patternIdx += 1;
    } else {
      sb += ch;
    }

    if (patternIdx == pattern.length && (idx+1) != compareString.length) {
      sb += str.substr(idx+1);
      break;
    }
  }

  if (patternIdx == pattern.length) {
    return {
      success: true,
      value: sb,
      score: calcScore(query, str, firstMatchIndex, lastMatchIndex-firstMatchIndex)
    };
  }

  return {
    success: false,
    score: 0
  };
}

function calcScore(query, str, firstIndex, matchLen) {
  var score = 100 * (query.length + 1) / ((1 + firstIndex) + (matchLen + 1));
  if (str.length-query.length < 5) {
    score = score + 20;
  } else if (str.length-query.length < 10) {
    score = score + 10;
  }

  return Math.floor(score);
}

module.exports = {
  match: match
};
