# Security Policy | 安全策略

## Supported Versions | 支持的版本

This project is primarily an educational resource for learning compression algorithms.
Security updates will be applied to the latest version on the `main` branch.

本项目主要是用于学习压缩算法的教育资源。安全更新将应用于 `main` 分支上的最新版本。

| Version | Supported          |
| ------- | ------------------ |
| main    | :white_check_mark: |
| others  | :x:                |

## Reporting a Vulnerability | 报告漏洞

We take security issues seriously, even in educational projects.

即使是教育项目，我们也认真对待安全问题。

### How to Report | 如何报告

**Preferred Method | 首选方式**: Use GitHub's private vulnerability reporting feature:

1. Go to the repository's **Security** tab
2. Click **Report a vulnerability**
3. Fill in the details

**Alternative | 备选方式**: If GitHub Security Advisories is not available, please open a regular issue with the title prefix `[SECURITY]`.

如果 GitHub Security Advisories 不可用，请提交一个标题前缀为 `[SECURITY]` 的普通 Issue。

### What to Include | 报告内容

Please include the following information in your report:

请在报告中包含以下信息：

- **Description | 描述**: A clear description of the vulnerability
- **Impact | 影响**: The potential impact of the vulnerability
- **Steps to Reproduce | 重现步骤**: Detailed steps to reproduce the issue
- **Affected Components | 受影响组件**: Which algorithm/language implementation is affected
- **Suggested Fix | 建议修复**: If you have a suggested fix, please include it

### Response Timeline | 响应时间

| Action | Timeline |
| ------ | -------- |
| Initial Response | Within 48 hours |
| Status Update | Within 7 days |
| Fix (if applicable) | Within 30 days |

| 操作 | 时间 |
| ---- | ---- |
| 初步响应 | 48 小时内 |
| 状态更新 | 7 天内 |
| 修复（如适用） | 30 天内 |

### What to Expect | 预期流程

1. **Acknowledgment | 确认**: We will acknowledge receipt of your report within 48 hours.
2. **Assessment | 评估**: We will assess the vulnerability and determine its severity.
3. **Communication | 沟通**: We will keep you informed of our progress.
4. **Resolution | 解决**: Once fixed, we will notify you and credit you (if desired) in the release notes.

## Scope | 范围

This security policy covers:

本安全策略涵盖：

- All algorithm implementations (Huffman, Arithmetic, Range coder, RLE)
- All language implementations (C++, Go, Rust)
- Benchmark and test scripts (Python)

### Out of Scope | 不在范围内

The following are generally out of scope:

以下通常不在范围内：

- Issues in third-party dependencies (please report to the respective projects)
- Issues that require physical access to the user's machine
- Social engineering attacks

## Security Best Practices | 安全最佳实践

When using this project:

使用本项目时：

1. **Validate Input | 验证输入**: Always validate input files before processing
2. **Resource Limits | 资源限制**: Be aware of memory usage when processing large files
3. **Trusted Sources | 可信来源**: Only decode files from trusted sources

## Acknowledgments | 致谢

We appreciate the security research community's efforts in helping keep this project safe.
Contributors who report valid security issues will be acknowledged in our release notes (unless they prefer to remain anonymous).

我们感谢安全研究社区为保持本项目安全所做的努力。报告有效安全问题的贡献者将在我们的发布说明中得到致谢（除非他们希望保持匿名）。
