document.addEventListener('DOMContentLoaded', function (index) {
    // 初始化 CodeMirror 编辑器
    const editor = CodeMirror.fromTextArea(document.getElementById('editorContent'), {
        mode: 'markdown',
        theme: 'default',
        lineNumbers: true,
        lineWrapping: false,
        autofocus: false,
        viewportMargin: Infinity,
        extraKeys: {
            "Enter": "newlineAndIndentContinueMarkdownList"
        },
        placeholder: "鼓励友善发言，禁止人身攻击"
    });

    const markedParse = marked
    markedParse.setOptions({
        gfm: true,
        tables: true,
        escaped : true,
        breaks: false,
        pedantic: false,
        sanitize: false,
        smartLists: true,
        smartypants: false,
    })
    console.log(markedParse, 'jjj---')

    // 设置编辑器高度
    editor.setSize('100%', '300px');

    // 获取预览相关元素
    const preview = document.getElementById('editorPreview');
    let htmlContent = ''




    // 后续再添加事件监听器
    // setTimeout(() => {
    //     editor.on('change', function(cm) {
    //         console.log('编辑器内容已更改');
    //         htmlContent = marked.parse(cm.getValue())
    //         console.log(htmlContent, 'oooo----')
    //     });
    // }, 0);



    // 表单提交
    const publishForm = document.getElementById('publishForm');
    publishForm.addEventListener('submit', function (e) {
        e.preventDefault();

        // 获取编辑器内容
        const markdownContent = editor.getValue();
        // const htmlContent = renderMarkdownToHtml(markdownContent);

        const formData = new FormData(this);
        const data = {
            title: formData.get('title'),
            category: formData.get('category'),
            tags: tags.join(','),
            content: htmlContent
        };

        // 发送到后端
        fetch('/api/posts', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        })
            .then(response => response.json())
            .then(result => {
                if (result.success) {
                    customAlert.success('文章发布成功！');
                    window.location.href = '/';
                } else {
                    customAlert.error('发布失败: ' + result.message);
                }
            })
            .catch(error => {
                console.error('Error:', error);
                customAlert.error('发布过程中出现错误');
            });
    });
});
