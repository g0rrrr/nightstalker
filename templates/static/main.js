$(document).ready(function() {
    const popups = {
        reply: {
            button: $('#top-action-button'),
            popup: $('#quickReplyPopup'),
            close: $('#reply_close'),
            header: $('#popupHeader')
        },
        edit: {
            button: $('#edit_button'),
            popup: $('#edit_popup'),
            close: $('#edit_close'),
            header: $('#edit-header')
        },
        post: {
            button: $('#post_button'),
            popup: $('#post_popup'),
            close: $('#post_close'),
            header: $('#post-header')
        }
    };

    $.each(popups, function(_, popupConfig) {
        const { button, popup, close, header } = popupConfig;

        button.on('click', function() {
            popup.css({ display: 'block', top: '100px', left: '100px' });
        });

        close.on('click', function() {
            popup.css('display', 'none');
        });

        let isDragging = false;
        let offsetX, offsetY;

        header.on('mousedown', function(e) {
            isDragging = true;
            offsetX = e.clientX - popup.offset().left;
            offsetY = e.clientY - popup.offset().top;
            $(document).on('mousemove', onMouseMove);
            $(document).on('mouseup', onMouseUp);
        });

        function onMouseMove(e) {
            if (isDragging) {
                popup.css({
                    left: `${e.clientX - offsetX}px`,
                    top: `${e.clientY - offsetY}px`
                });
            }
        }

        function onMouseUp() {
            isDragging = false;
            $(document).off('mousemove', onMouseMove);
            $(document).off('mouseup', onMouseUp);
        }
    });

    // Original Sirjtaa functionality
    function quotePostClicked(e) {
        e.preventDefault();

        var pid = $(this).data("postid");
        var content = $("#p" + pid + "-unparsed-content").text();
        var author  = $("#p" + pid + "-author").text();

        // Remove any previous quotes from the content
        content = content.replace(/.* said:\s*(\>.*\n)*/g, "");
        content = content.replace(/>.*\n/g, "").trim();
        content = content.replace(/\n/g, "\n>");

        var reply = $("#reply-field");
        var val  = ">" + author + "\n>" + content + "\n";
        reply.focus().val('').val(val).scrollTop(reply[0].scrollHeight);
    }

    function threadReplyClicked() {
        $("#reply-field").focus();
    }

    function clickConfirmDelete() {
        return confirm("Really delete this? It can't be undone.");
    }

    function clickModerate(e) {
        e.preventDefault();
        $(this).hide();
        $(this).siblings(".mod_tools").show();
        $(this).siblings(".prec_slash").hide();
    }

    $(".quote-post").on("click", quotePostClicked);
    $("#top-action-button").on("click", threadReplyClicked);
    $(".delete").on("click", clickConfirmDelete);
    $(".moderate").on("click", clickModerate);
});

