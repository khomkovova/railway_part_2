



new Vue({
    el: '#app',
    data:{
        error:"",
        status:"",
        comments:""
    },
    methods:{

        submit(){
            axios({ method: "POST", "url": "http:api/addcomments", data:'{"comments":'+ '"' + this.comments + '"}'}).then(result => {
                this.status = "Your comment send";
            }, error => {
                this.status = "Bad username or password";
            });

        }
    }
});