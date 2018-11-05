
new Vue({
    el: '#main',
    data:{
        error:"",
        status:"",
        allComments:""
    },
    mounted(){
        axios({ method: "GET", "url": "http:/api/getcomments" }).then(result => {
            console.error(result.data);
            this.allComments = result.data["comments"];
        }, error => {
            console.error(error);
        });
    },
    methods:{

        getComments() {
            axios({ method: "GET", "url": "http:/api/getcomments" }).then(result => {
                console.error(result.data);
                this.allComments = result.data["comments"];
            }, error => {
                console.error(error);
            });
        },

        addComment(){
            axios({ method: "POST", "url": "http:/api/setcomments","data":"test comments" }).then(result => {
                console.error(result.data);

            }, error => {
                console.error(error);
            });
        }

    }
});