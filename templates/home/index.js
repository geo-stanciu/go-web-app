function f1() {
    postAJAX('/exchange-rates', { test1: 'a', test2: 'b' }, function(data) {
        console.log('1', data);
    });

    getAJAX('/exchange-rates', { test1: 'a', test2: 'b' }, function(data) {
        console.log('2', data);
    });

    return "1";
}