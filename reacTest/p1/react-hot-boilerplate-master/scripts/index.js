'use strict';

var React = require('react'),
    CommentBox = require('./app');

React.render(<CommentBox url="comments.json"/>, document.body);
