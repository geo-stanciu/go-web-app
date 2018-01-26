function setDstDetails() {
    var now = new Date();
    var later = new Date();
    // Set time for how long the cookie should be saved
    later.setTime(now.getTime() + 365 * 24 * 60 * 60 * 1000);

    // Set cookie for the time zone offset in minutes
    Cookies.set("time_zone_offset", (-now.getTimezoneOffset()).toString(), { expires: later, path: "/" });

    Cookies.getTime

    // Create two new dates
    var d1 = new Date();
    var d2 = new Date();

    // Date one is set to January 1st of this year
    // Guaranteed not to be in DST for northern hemisphere,
    // and guaranteed to be in DST for southern hemisphere
    // (If DST exists on client PC)
    d1.setDate(1);
    d1.setMonth(1);

    // Date two is set to July 1st of this year
    // Guaranteed to be in DST for northern hemisphere,
    // and guaranteed not to be in DST for southern hemisphere
    // (If DST exists on client PC)
    d2.setDate(1);
    d2.setMonth(7);

    // If time zone offsets match, no DST exists for this time zone
    if (parseInt(d1.getTimezoneOffset()) == parseInt(d2.getTimezoneOffset())) {
        Cookies.set("time_zone_dst", "0", { expires: later, path: "/" });
    }
    // DST exists for this time zone – check if it is currently active
    else {
        // Find out if we are on northern or southern hemisphere
        // Hemisphere is positive for northern, and negative for southern
        var hemisphere = parseInt(d1.getTimezoneOffset()) - parseInt(d2.getTimezoneOffset());

        // Current date is still before or after DST, not containing DST
        if ((hemisphere > 0 && parseInt(d1.getTimezoneOffset()) == parseInt(now.getTimezoneOffset())) ||
            (hemisphere < 0 && parseInt(d2.getTimezoneOffset()) == parseInt(now.getTimezoneOffset()))) {
            Cookies.set("time_zone_dst", "0", { expires: later, path: "/" });
        } // DST is active right now with the current date
        else {
            Cookies.set("time_zone_dst", "1", { expires: later, path: "/" });
        }
    }
}

setDstDetails();