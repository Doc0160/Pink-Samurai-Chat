"use strict";

var chat = {
    "msg": "",
    "username": "",
    "channel": "",
    "log": "",
    
    "conn": "",

    "message_senders": {
        "SendMessage": function(){
            if (!chat.conn) {
			    return;
		    }
	        chat.conn.send(JSON.stringify({
			    "type":     "message",
			    "channel":  chat.channel.value,
			    "text":     chat.msg.value
		    }));
        },
        "SendCommand": function() {
            if (!chat.conn) {
			    return;
		    }
            chat.conn.send(JSON.stringify({
			    "type":     "command",
			    "channel":  chat.channel.value,
			    "command":  chat.msg.value
		    }));
        },
    },
    
    "message_handlers": {
	    "disconnect": function(item, decoded){
            item.innerText = decoded.username + " got disconnected"
	        item.setAttribute("class", item.getAttribute("class") + " dischat.connect");
        },
	    "message": function(item, decoded){
            item.innerText = decoded.username + " : " + decoded.text;
	        if (!document.hasFocus())
		        new Notification(decoded.username, {
		            "body": decoded.text.replace("\\n", "\n"),
		        });
	        item.innerHTML = chat.convert_urls(item.innerHTML);
	        item.innerHTML = emojione.toImage(item.innerHTML);
	        item.setAttribute("class", item.getAttribute("class") + " message");
        },
	    "channel_leave": function(item, decoded){
            item.innerText = decoded.username + " quited " + decoded.channel;
	        item.setAttribute("class", item.getAttribute("class") + " channel_quit");
        },
	    "channel_join": function(item, decoded){
            item.innerText = decoded.username + " joined " + decoded.channel;
	        item.setAttribute("class", item.getAttribute("class") + " channel_join");
        },
    },
    
    "getUsername": function(){
        return chat.getCookie("username");
    },
    
    "getCookie": function(cname) {
        var name = cname + "=";
        var ca = document.cookie.split(';');
        for(var i = 0; i <ca.length; i++) {
            var c = ca[i];
            while (c.charAt(0)==' ') {
                c = c.substring(1);
            }
            if (c.indexOf(name) == 0) {
                return c.substring(name.length,c.length);
            }
        }
        return "";
    },

    "randomName": function(){
        var names = ['Dick',      'Fag',    'Pussy', 'Shit',
			         'HorseCock', 'Nigger', 'Dog',   'Ball',
			         'Cock',      'Anus',   'Cat',   'Cum',
			         'Gay',       'Banana', 'Squid', 'Skull',
			         'Slut',      'Ass',    'Vagina','Alien',
			         'JustinBieber', 'Tentacule-Sama'];
        var actions_ing = ['Farting',    'Baiting',   'Shitting',
				           'Licking',    'Swallowing', 'Consuming',
				           'Sucking'];
        var actions_er = ['Sucker',    'Licker',   'Snorter', 'Warmer',
			              'Fucker',    'Slurper',  'Wanker',  'Biter',
			              'Destroyer', 'Worker',   'Swallower',
			              'Lover',     'Massager', 'Researcher'];
        return names[Math.floor(Math.random()*names.length)]
		      +actions_ing[Math.floor(Math.random()*actions_ing.length)]
		      +" "
		      +names[Math.floor(Math.random()*names.length)]
		      +actions_er[Math.floor(Math.random()*actions_er.length)];
	},
    
    "convert_urls": function(text){
		text = text.replace("\\n", "<br/>");
		var mime_types = {
		    mp3: "audio/mp3",
		    mp4: "video/mp4",
		    webm: "video/webm",
		    png: "image/png",
		    gif: "image/gif",
		    jpg: "image/jpeg",
		    jpeg: "image/jpeg",
		    jpe: "image/jpeg",
		    svf: "image/svg+xml",
		};
		var exp = /(\b(\w+):(\/+)((\w|[%\-\.\/\?=&])+))/ig;
		text = text.replace(exp, function(match){
		    var ext = match.split('.').pop().toLowerCase();
		    if(mime_types[ext] != undefined) {
			    var mime = mime_types[ext];
			    var type = mime.split('/').shift();
			    if(type == "image") {
			        return "<a href='"+match+"'><img src='"+match+"'/></a>";
			    } else if(type == "video") {
			        return "<video controls><source src='"+match+"' type='"+mime+"'>Your browser does not support the video element.</video>";
			    } else if(type == "audio") {
			        return "<audio controls><source src='"+match+"' type='"+mime+"'>Your browser does not support the audio element.</audio>";
			    }
		    } else {
			    return "<a href='"+match+"'>"+match+"</a>";
		    }
		});
		var exp = /(?:([\*_~]{1,3}))([^\*_~\n]+[^\*_~\s])\1/gi
		var text = text.replace(exp, function(match, one, two){
		    if(one == "~~"){
			    return '<del>'+two+'</del>';
		    }else if(one.length == 1){
			    return '<b>'+two+'</b>';
		    }else if(one.length == 2){
			    return '<i>'+two+'</i>';
		    }else if(one.length == 3){
			    return '<b><i>'+two+'</i></b>';
		    }else{
			    return two;
		    }
		});
		var exp = /\&gt;(.+)/
		var text = text.replace(exp, '<code>$1</code>');
		var exp = /\s\`\`\`\n?([^`]+)\`\`\`/g
		var text = text.replace(exp, '<code>$1</code>');
		return text;
    }
};
