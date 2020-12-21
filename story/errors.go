package story

import "errors"

var ErrDataTooBig = errors.New("data too big")
var ErrHeaderMissing = errors.New("header missing")
var ErrInvalidStory = errors.New("data is not a valid story")
