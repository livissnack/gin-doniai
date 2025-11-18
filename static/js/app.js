// 自定义Alert组件
class CustomAlert {
  constructor() {
    this.container = document.createElement('div');
    this.container.id = 'custom-alert-container';
    this.container.style.position = 'fixed';
    this.container.style.top = '0';
    this.container.style.right = '0';
    this.container.style.zIndex = '10000';
    this.container.style.width = '100%';
    this.container.style.maxWidth = '400px';
    this.container.style.padding = '20px';
    document.body.appendChild(this.container);
  }

  show(message, type = 'info', duration = 3000) {
    // 创建alert元素
    const alert = document.createElement('div');
    alert.className = `custom-alert ${type}`;

    // 根据类型设置图标
    let icon = 'ℹ️';
    if (type === 'success') icon = '✅';
    else if (type === 'error') icon = '❌';
    else if (type === 'warning') icon = '⚠️';

    alert.innerHTML = `
      <div class="alert-content">
        <span class="alert-icon">${icon}</span>
        <span class="alert-message">${message}</span>
        <button class="alert-close">&times;</button>
      </div>
    `;

    // 添加关闭事件
    const closeBtn = alert.querySelector('.alert-close');
    closeBtn.addEventListener('click', () => {
      this.hide(alert);
    });

    // 添加到容器
    this.container.appendChild(alert);

    // 触发显示动画
    setTimeout(() => {
      alert.classList.add('show');
    }, 10);

    // 自动关闭
    if (duration > 0) {
      setTimeout(() => {
        this.hide(alert);
      }, duration);
    }

    return alert;
  }

  hide(alert) {
    alert.classList.remove('show');
    setTimeout(() => {
      if (alert.parentNode) {
        alert.parentNode.removeChild(alert);
      }
    }, 300);
  }

  success(message, duration) {
    return this.show(message, 'success', duration);
  }

  error(message, duration) {
    return this.show(message, 'error', duration);
  }

  warning(message, duration) {
    return this.show(message, 'warning', duration);
  }

  info(message, duration) {
    return this.show(message, 'info', duration);
  }
}

// 创建全局实例
const customAlert = new CustomAlert();

// 主题切换功能
class ThemeManager {
  constructor() {
    this.themeToggle = document.getElementById('themeToggle');
    this.themeIcon = this.themeToggle.querySelector('.theme-icon');
    this.body = document.body;

    this.init();
  }

  init() {
    // 从localStorage加载用户主题偏好
    const savedTheme = localStorage.getItem('theme') || 'dark-theme';
    this.setTheme(savedTheme);

    // 绑定切换事件
    this.themeToggle.addEventListener('click', () => this.toggleTheme());

    // 添加主题切换动画类
    this.body.classList.add('theme-transition');
  }

  toggleTheme() {
    const isDark = this.body.classList.contains('dark-theme');
    const newTheme = isDark ? 'light-theme' : 'dark-theme';
    this.setTheme(newTheme);
  }

  setTheme(theme) {
    // 移除现有主题类
    this.body.classList.remove('dark-theme', 'light-theme');

    // 添加新主题类
    this.body.classList.add(theme);

    // 更新图标
    this.updateIcon(theme);

    // 保存到localStorage
    localStorage.setItem('theme', theme);

    // 触发自定义事件（便于其他组件监听主题变化）
    window.dispatchEvent(new CustomEvent('themeChanged', { detail: theme }));
  }

  updateIcon(theme) {
    const isDark = theme === 'dark-theme';
    this.themeIcon.textContent = isDark ? '🌙' : '☀️';
    this.themeToggle.setAttribute('title', isDark ? '切换到亮色主题' : '切换到暗黑主题');
  }

  // 获取当前主题
  getCurrentTheme() {
    return this.body.classList.contains('dark-theme') ? 'dark-theme' : 'light-theme';
  }
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
  window.themeManagerInstance = new ThemeManager(); // 保存实例

  window.addEventListener('themeChanged', function(event) { // 使用普通函数
    console.log('主题已切换至1:', event.detail);
    // window.themeManagerInstance.toggleTheme();
    console.log('主题已切换至2:', event.detail);
  });

  // 添加一些交互效果
});


// 用户下拉菜单功能
document.addEventListener('DOMContentLoaded', function() {
  const userDropdown = document.getElementById('userDropdown');
  const dropdownMenu = document.getElementById('dropdownMenu');

  if (userDropdown && dropdownMenu) {
    userDropdown.addEventListener('click', function(e) {
      e.stopPropagation();
      dropdownMenu.style.display = dropdownMenu.style.display === 'block' ? 'none' : 'block';
    });

    // 点击其他地方关闭下拉菜单
    document.addEventListener('click', function() {
      dropdownMenu.style.display = 'none';
    });
  }
});


// 设置页面功能
document.addEventListener('DOMContentLoaded', function() {
  // 个人资料表单提交
  const profileForm = document.getElementById('profileForm');
  if (profileForm) {
    profileForm.addEventListener('submit', function(e) {
      e.preventDefault();
      // 这里可以添加保存个人资料的逻辑
      // alert('个人资料保存成功！');
      customAlert.info('个人资料保存成功');
    });
  }

  // 安全设置表单提交
  const securityForm = document.getElementById('securityForm');
  if (securityForm) {
    securityForm.addEventListener('submit', function(e) {
      e.preventDefault();
      // 这里可以添加更改密码的逻辑
      // alert('密码更改成功！');
      customAlert.info('密码更改成功！');
    });
  }
});

// 滚动到顶部/底部功能
class ScrollManager {
  constructor() {
    this.scrollTopBtn = document.getElementById('scrollTopBtn');
    this.scrollBottomBtn = document.getElementById('scrollBottomBtn');
    this.scrollThreshold = 300; // 滚动超过300px时显示按钮

    this.init();
  }

  init() {
    if (this.scrollTopBtn && this.scrollBottomBtn) {
      // 绑定滚动事件
      window.addEventListener('scroll', () => this.handleScroll());

      // 绑定按钮点击事件
      this.scrollTopBtn.addEventListener('click', () => this.scrollToTop());
      this.scrollBottomBtn.addEventListener('click', () => this.scrollToBottom());

      // 初始检查
      this.handleScroll();
    }
  }

  handleScroll() {
    // 检查滚动位置
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    const scrollHeight = document.documentElement.scrollHeight;
    const clientHeight = document.documentElement.clientHeight;
    const scrolledToBottom = scrollTop + clientHeight >= scrollHeight - 5;

    // 显示/隐藏回到顶部按钮
    if (scrollTop > this.scrollThreshold) {
      this.scrollTopBtn.classList.remove('hidden');
    } else {
      this.scrollTopBtn.classList.add('hidden');
    }

    // 显示/隐藏回到底部按钮
    if (!scrolledToBottom) {
      this.scrollBottomBtn.classList.remove('hidden');
    } else {
      this.scrollBottomBtn.classList.add('hidden');
    }
  }

  scrollToTop() {
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    });
  }

  scrollToBottom() {
    window.scrollTo({
      top: document.documentElement.scrollHeight,
      behavior: 'smooth'
    });
  }
}

// 页面加载完成后初始化滚动管理器
document.addEventListener('DOMContentLoaded', () => {
  // 初始化滚动管理器
  window.scrollManager = new ScrollManager();
});



