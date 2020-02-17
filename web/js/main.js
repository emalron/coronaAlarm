const app = new Vue({
    el: '#root',
    data: {
        items: [],
        input: "",
        result: "",
    },
    methods: {
        test: function(e) {
            e.preventDefault();
            let inputs = this.input.split(' ').slice(0,2) 
            getData(this, inputs, function(obj, data, num) {
                let result = JSON.parse(data)
                for(var i=0; i<num; i++) {
                    result = highlight(result, inputs[i])
                }
                obj.input = "";
                obj.items = result
                obj.$refs.answer.focus();
            })
        },
    },
});

var getData = function(obj, search, callback) {
    let xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if(xhr.readyState === xhr.DONE) {
            if(xhr.status === 200 || xhr.status === 201) {
                callback(obj, xhr.responseText, search.length)
            }
        }
    }
    let url = "http://api.emalron.com:8001/v1/search";
    switch(search.length) {
        case 0:
            break;
        case 1:
            url += "/" + search[0]
            break;
        case 2:
            url += "/" + search[0] + "/" + search[1]
    }

    xhr.open('GET', url)
    xhr.send()
}

var highlight = function(data, search) {
    console.log(`data: ${data}, num: ${data.length}`)
    let num = data.length
    for (var i=0; i<num; i++) {
        let route = data[i].route
        let nroute = route.length
        for (var j=0; j<nroute; j++) {
            route[j] = route[j].split(search).join(`<span class="highlight">${search}</span>`)
            console.log(route[j])
        }
    }
    return data
}