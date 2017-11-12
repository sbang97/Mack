"use strict"

var searchForm = document.querySelector(".search-form");
var searchInput = searchForm.querySelector("input");
var img = document.getElementById("img");
var title = document.getElementById("title");
var description = document.getElementById("description");


searchForm.addEventListener("submit", function(evt) {
    evt.preventDefault();
    var query = searchInput.value.trim();

    if (query.length <= 0) {
        return false;
    }
    fetch("https://api.sbang9.me/v1/summary?url=" + query)
        .then(function(resp) {
            return resp.json()
        })
        .then(function(data) {
            if (data["title"] != null ) {
                title.textContent = data["title"];
            } else {
                title.textContent = "No open graph title found...";
            }
            if (data["description"] != null ) {
                description.textContent = data["description"]
            } else {
                description.textContent = "No open graph description found...";
            }
            if (data["image"] != null && data["image"] != "") {
                img.src = data["image"];
            } else {
                img.alt = "Image not found...";
                img.src= "";
            }
        })
        .catch(function(error) {
            alert(error);
        })
})