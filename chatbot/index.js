"use strict";

const express = require('express');
const { Wit } = require('node-wit');
const morgan = require('morgan');
var bodyParser = require('body-parser');

const app = express();
const port = process.env.PORT || '80';
const host = process.env.HOST || '';
const witaiToken = process.env.WITAITOKEN;

if (!witaiToken) {
	console.error("please set WITAITOKEN to your wit.ai app token");
	process.exit(1);
}

const witaiClient = new Wit({ accessToken: witaiToken });

var MongoClient = require('mongodb').MongoClient;
var assert = require('assert');

// Connection URL
var url = 'mongodb://api.sbang9.me:27017';

function handleLastPost(req, res, witaiData) {
	// Use connect method to connect to the server
	MongoClient.connect(url, function(err, db) {
		assert.equal(null, err);

		var collection = db.collection('messages');
		// Find some documents
		collection.find({}).sort({'datetime': -1}).limit(1).toArray(function(err, docs) {
			assert.equal(err, null);
			res.send(docs);
		});
		db.close();
	});
}

function handleNumPosts(req, res, witaiData) {
		// Use connect method to connect to the server
	MongoClient.connect(url, function(err, db) {
		assert.equal(null, err);

		var collection = db.collection('messages');
		// Find some documents
		collection.count({'channel': witAiData}).toArray(function(err, docs) {
			assert.equal(err, null);
			res.send(docs);
		});
		db.close();
	});
}

function handleUsersinChannel(req, res, witaiData) {
	// Use connect method to connect to the server
	MongoClient.connect(url, function(err, db) {
		assert.equal(null, err);

		var collection = db.collection('messages');
		// Find some documents
		collection.find({'channel': witAiData}).toArray(function(err, docs) {
			assert.equal(err, null);
			res.send(docs);
		});
		db.close();
	});
}

function handelHavePosted(req, res, witaiData) {
	// Use connect method to connect to the server
	MongoClient.connect(url, function(err, db) {
		assert.equal(null, err);

		var collection = db.collection('messages');
		// Find some documents
		collection.find({}).toArray(function(err, docs) {
			assert.equal(err, null);
			res.send(docs);
		});
		db.close();
	});
}

function handelHaveNotPosted(req, res, witaiData) {
	// Use connect method to connect to the server
	MongoClient.connect(url, function(err, db) {
		assert.equal(null, err);

		var collection = db.collection('messages');
		// Find some documents
		collection.find({}).toArray(function(err, docs) {
			assert.equal(err, null);
			res.send(docs);
		});
		db.close();
	});
}
// ask how to query the database and get the answers.
app.get("/v1/bot", (req, res, next) => {
    if (req.method.toLowerCase() != "post") {
		res.send("the request method must be post");
		break;
	}
	let query = bodyParser.text(req.body);
    witaiClient.message(query)
        .then(data => {
            switch (data.entities.intent[0].value) {
				case "last":
					handleLastPost(req, res, data);
					break;
				case "num-posts":
					handleNumPosts(req, res, data);
					break;
				case "users":
					handleUsersinChannel(req, res, data);
					break;
				case "have-posted":
					handleHavePosted(req, res, data);
					break;
				case "not-posted":
					handleHaveNotPosted(req, res, data);
					break;
				default:
					res.send("Sorry, im not sure how to answer that :(");
			}
        })
        .catch(next); 
});