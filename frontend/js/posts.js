$(document).ready(function() {
    // Function to fetch and display posts
    function fetchPosts() {
        console.log('Fetching posts...'); // Debugging
        $.ajax({
            url: 'http://localhost:8080/posts',
            method: 'GET',
            success: function(posts) {
                console.log('Posts fetched:', posts); // Debugging
                $('#postsList').empty();
                posts.forEach(post => {
                    const comments = post.comments || []; // Ensure comments is an array
                    $('#postsList').append(`
                        <div class="card mb-3">
                            <div class="card-body">
                                <p class="card-text">${post.content}</p>
                                <button class="btn btn-success like-btn" data-id="${post.id}">Like (${post.likes})</button>
                                <button class="btn btn-danger dislike-btn" data-id="${post.id}">Dislike (${post.dislikes})</button>
                                <button class="btn btn-secondary share-btn" data-id="${post.id}">Share</button>
                                <div class="mt-2">
                                    <textarea class="form-control comment-text" rows="2" placeholder="Add a comment"></textarea>
                                    <button class="btn btn-primary mt-2 comment-btn" data-id="${post.id}">Comment</button>
                                </div>
                            </div>
                            <div class="card-footer">
                                <div class="comments-list">
                                    ${comments.map(comment => `<p>${comment.content}</p>`).join('')}
                                </div>
                            </div>
                        </div>
                    `);
                });
            },
            error: function(xhr, status, error) {
                console.error('Error fetching posts:', xhr.responseText);
            }
        });
    }

    // Fetch posts on page load
    fetchPosts();

    // Handle post form submission
    $('#postForm').submit(function(event) {
        event.preventDefault();
        const content = $('#postContent').val();
        console.log('Posting new content:', content); // Debugging
        $.ajax({
            url: 'http://localhost:8080/add_post',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ content: content }),
            success: function() {
                $('#postContent').val('');
                fetchPosts();
            },
            error: function(xhr, status, error) {
                console.error('Error adding post:', xhr.responseText);
            }
        });
    });

    // Handle like button click
    $(document).on('click', '.like-btn', function() {
        const postId = $(this).data('id');
        console.log('Liking post:', postId); // Debugging
        $.ajax({
            url: 'http://localhost:8080/like_post',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ post_id: postId }),
            success: function() {
                fetchPosts();
            },
            error: function(xhr, status, error) {
                console.error('Error liking post:', xhr.responseText);
            }
        });
    });

    // Handle dislike button click
    $(document).on('click', '.dislike-btn', function() {
        const postId = $(this).data('id');
        console.log('Disliking post:', postId); // Debugging
        $.ajax({
            url: 'http://localhost:8080/dislike_post',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ post_id: postId }),
            success: function() {
                fetchPosts();
            },
            error: function(xhr, status, error) {
                console.error('Error disliking post:', xhr.responseText);
            }
        });
    });

    // Handle share button click
    $(document).on('click', '.share-btn', function() {
        const postId = $(this).data('id');
        console.log('Sharing post:', postId); // Debugging
        $.ajax({
            url: 'http://localhost:8080/share_post',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ post_id: postId }),
            success: function(response) {
                alert(`Share link: ${response.share_link}`);
            },
            error: function(xhr, status, error) {
                console.error('Error sharing post:', xhr.responseText);
            }
        });
    });

    // Handle comment button click
    $(document).on('click', '.comment-btn', function() {
        const postId = $(this).data('id');
        const content = $(this).closest('.card').find('.comment-text').val();
        console.log('Commenting on post:', postId, 'Content:', content); // Debugging
        $.ajax({
            url: 'http://localhost:8080/add_comment',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ post_id: postId, content: content }),
            success: function() {
                fetchPosts();
            },
            error: function(xhr, status, error) {
                console.error('Error adding comment:', xhr.responseText);
            }
        });
    });
});
