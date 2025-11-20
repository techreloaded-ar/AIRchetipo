const path = require('node:path');
const { BaseIdeSetup } = require('./_base-ide');
const chalk = require('chalk');
const { getAgentsFromAir, getTasksFromAir } = require('./shared/air-artifacts');
const { AgentCommandGenerator } = require('./shared/agent-command-generator');

/**
 * Qwen Code setup handler
 * Creates TOML command files in .qwen/commands/AIRchetipo/
 */
class QwenSetup extends BaseIdeSetup {
  constructor() {
    super('qwen', 'Qwen Code');
    this.configDir = '.qwen';
    this.commandsDir = 'commands';
    this.airDir = 'air';
  }

  /**
   * Setup Qwen Code configuration
   * @param {string} projectDir - Project directory
   * @param {string} airDir - AIRCHETIPO installation directory
   * @param {Object} options - Setup options
   */
  async setup(projectDir, airDir, options = {}) {
    console.log(chalk.cyan(`Setting up ${this.name}...`));

    // Create .qwen/commands/AIRchetipo directory structure
    const qwenDir = path.join(projectDir, this.configDir);
    const commandsDir = path.join(qwenDir, this.commandsDir);
    const airCommandsDir = path.join(commandsDir, this.airDir);

    await this.ensureDir(airCommandsDir);

    // Update existing settings.json if present
    await this.updateSettings(qwenDir);

    // Clean up old configuration if exists
    await this.cleanupOldConfig(qwenDir);

    // Generate agent launchers
    const agentGen = new AgentCommandGenerator(this.airFolderName);
    const { artifacts: agentArtifacts } = await agentGen.collectAgentArtifacts(airDir, options.selectedModules || []);

    // Get tasks, tools, and workflows (standalone only for tools/workflows)
    const tasks = await getTasksFromAir(airDir, options.selectedModules || []);
    const tools = await this.getTools(airDir, true);
    const workflows = await this.getWorkflows(airDir, true);

    // Create directories for each module (including standalone)
    const modules = new Set();
    for (const item of [...agentArtifacts, ...tasks, ...tools, ...workflows]) modules.add(item.module);

    for (const module of modules) {
      await this.ensureDir(path.join(airCommandsDir, module));
      await this.ensureDir(path.join(airCommandsDir, module, 'agents'));
      await this.ensureDir(path.join(airCommandsDir, module, 'tasks'));
      await this.ensureDir(path.join(airCommandsDir, module, 'tools'));
      await this.ensureDir(path.join(airCommandsDir, module, 'workflows'));
    }

    // Create TOML files for each agent launcher
    let agentCount = 0;
    for (const artifact of agentArtifacts) {
      // Convert markdown launcher content to TOML format
      const tomlContent = this.processAgentLauncherContent(artifact.content, {
        module: artifact.module,
        name: artifact.name,
      });

      const targetPath = path.join(airCommandsDir, artifact.module, 'agents', `${artifact.name}.toml`);

      await this.writeFile(targetPath, tomlContent);

      agentCount++;
      console.log(chalk.green(`  ✓ Added agent: /air:${artifact.module}:agents:${artifact.name}`));
    }

    // Create TOML files for each task
    let taskCount = 0;
    for (const task of tasks) {
      const content = await this.readAndProcess(task.path, {
        module: task.module,
        name: task.name,
      });

      const targetPath = path.join(airCommandsDir, task.module, 'tasks', `${task.name}.toml`);

      await this.writeFile(targetPath, content);

      taskCount++;
      console.log(chalk.green(`  ✓ Added task: /air:${task.module}:tasks:${task.name}`));
    }

    // Create TOML files for each tool
    let toolCount = 0;
    for (const tool of tools) {
      const content = await this.readAndProcess(tool.path, {
        module: tool.module,
        name: tool.name,
      });

      const targetPath = path.join(airCommandsDir, tool.module, 'tools', `${tool.name}.toml`);

      await this.writeFile(targetPath, content);

      toolCount++;
      console.log(chalk.green(`  ✓ Added tool: /air:${tool.module}:tools:${tool.name}`));
    }

    // Create TOML files for each workflow
    let workflowCount = 0;
    for (const workflow of workflows) {
      const content = await this.readAndProcess(workflow.path, {
        module: workflow.module,
        name: workflow.name,
      });

      const targetPath = path.join(airCommandsDir, workflow.module, 'workflows', `${workflow.name}.toml`);

      await this.writeFile(targetPath, content);

      workflowCount++;
      console.log(chalk.green(`  ✓ Added workflow: /air:${workflow.module}:workflows:${workflow.name}`));
    }

    console.log(chalk.green(`✓ ${this.name} configured:`));
    console.log(chalk.dim(`  - ${agentCount} agents configured`));
    console.log(chalk.dim(`  - ${taskCount} tasks configured`));
    console.log(chalk.dim(`  - ${toolCount} tools configured`));
    console.log(chalk.dim(`  - ${workflowCount} workflows configured`));
    console.log(chalk.dim(`  - Commands directory: ${path.relative(projectDir, airCommandsDir)}`));

    return {
      success: true,
      agents: agentCount,
      tasks: taskCount,
      tools: toolCount,
      workflows: workflowCount,
    };
  }

  /**
   * Update settings.json to remove old agent references
   */
  async updateSettings(qwenDir) {
    const fs = require('fs-extra');
    const settingsPath = path.join(qwenDir, 'settings.json');

    if (await fs.pathExists(settingsPath)) {
      try {
        const settingsContent = await fs.readFile(settingsPath, 'utf8');
        const settings = JSON.parse(settingsContent);
        let updated = false;

        // Remove agent file references from contextFileName
        if (settings.contextFileName && Array.isArray(settings.contextFileName)) {
          const originalLength = settings.contextFileName.length;
          settings.contextFileName = settings.contextFileName.filter(
            (fileName) => !fileName.startsWith('agents/') && !fileName.startsWith('airchetipo/'),
          );

          if (settings.contextFileName.length !== originalLength) {
            updated = true;
          }
        }

        if (updated) {
          await fs.writeFile(settingsPath, JSON.stringify(settings, null, 2));
          console.log(chalk.green('  ✓ Updated .qwen/settings.json'));
        }
      } catch (error) {
        console.warn(chalk.yellow('  ⚠ Could not update settings.json:'), error.message);
      }
    }
  }

  /**
   * Clean up old configuration directories
   */
  async cleanupOldConfig(qwenDir) {
    const fs = require('fs-extra');
    const agentsDir = path.join(qwenDir, 'agents');
    const airMethodDir = path.join(qwenDir, 'airchetipo');
    const airDir = path.join(qwenDir, 'airDir');

    if (await fs.pathExists(agentsDir)) {
      await fs.remove(agentsDir);
      console.log(chalk.green('  ✓ Removed old agents directory'));
    }

    if (await fs.pathExists(airMethodDir)) {
      await fs.remove(airMethodDir);
      console.log(chalk.green('  ✓ Removed old airchetipo directory'));
    }

    if (await fs.pathExists(airDir)) {
      await fs.remove(airDir);
      console.log(chalk.green('  ✓ Removed old AIRchetipo directory'));
    }
  }

  /**
   * Read and process file content
   */
  async readAndProcess(filePath, metadata) {
    const fs = require('fs-extra');
    const content = await fs.readFile(filePath, 'utf8');
    return this.processContent(content, metadata);
  }

  /**
   * Process agent launcher content and convert to TOML format
   * @param {string} launcherContent - Launcher markdown content
   * @param {Object} metadata - File metadata
   * @returns {string} TOML formatted content
   */
  processAgentLauncherContent(launcherContent, metadata = {}) {
    // Strip frontmatter from launcher content
    const frontmatterRegex = /^---\s*\n[\s\S]*?\n---\s*\n/;
    const contentWithoutFrontmatter = launcherContent.replace(frontmatterRegex, '');

    // Extract title for TOML description
    const titleMatch = launcherContent.match(/description:\s*"([^"]+)"/);
    const title = titleMatch ? titleMatch[1] : metadata.name;

    // Create TOML with launcher content (without frontmatter)
    return `description = "AIRCHETIPO ${metadata.module.toUpperCase()} Agent: ${title}"
prompt = """
${contentWithoutFrontmatter.trim()}
"""
`;
  }

  /**
   * Override processContent to add TOML metadata header for Qwen
   * @param {string} content - File content
   * @param {Object} metadata - File metadata
   * @returns {string} Processed content with Qwen template
   */
  processContent(content, metadata = {}) {
    // First apply base processing (includes activation injection for agents)
    let prompt = super.processContent(content, metadata);

    // Determine the type and description based on content
    const isAgent = content.includes('<agent');
    const isTask = content.includes('<task');
    const isTool = content.includes('<tool');
    const isWorkflow = content.includes('workflow:') || content.includes('name:');

    let description = '';

    if (isAgent) {
      // Extract agent title if available
      const titleMatch = content.match(/title="([^"]+)"/);
      const title = titleMatch ? titleMatch[1] : metadata.name;
      description = `AIRCHETIPO ${metadata.module.toUpperCase()} Agent: ${title}`;
    } else if (isTask) {
      // Extract task name if available
      const nameMatch = content.match(/name="([^"]+)"/);
      const taskName = nameMatch ? nameMatch[1] : metadata.name;
      description = `AIRCHETIPO ${metadata.module.toUpperCase()} Task: ${taskName}`;
    } else if (isTool) {
      // Extract tool name if available
      const nameMatch = content.match(/name="([^"]+)"/);
      const toolName = nameMatch ? nameMatch[1] : metadata.name;
      description = `AIRCHETIPO ${metadata.module.toUpperCase()} Tool: ${toolName}`;
    } else if (isWorkflow) {
      // Workflow
      description = `AIRCHETIPO ${metadata.module.toUpperCase()} Workflow: ${metadata.name}`;
    } else {
      description = `AIRCHETIPO ${metadata.module.toUpperCase()}: ${metadata.name}`;
    }

    return `description = "${description}"
prompt = """
${prompt}
"""
`;
  }

  /**
   * Format name as title
   */
  formatTitle(name) {
    return name
      .split('-')
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ');
  }

  /**
   * Cleanup Qwen configuration
   */
  async cleanup(projectDir) {
    const fs = require('fs-extra');
    const airCommandsDir = path.join(projectDir, this.configDir, this.commandsDir, this.airDir);
    const oldAirMethodDir = path.join(projectDir, this.configDir, 'airchetipo');
    const oldAIRchetipoDir = path.join(projectDir, this.configDir, 'AIRchetipo');

    if (await fs.pathExists(airCommandsDir)) {
      await fs.remove(airCommandsDir);
      console.log(chalk.dim(`Removed AIRchetipo configuration from Qwen Code`));
    }

    if (await fs.pathExists(oldAirMethodDir)) {
      await fs.remove(oldAirMethodDir);
      console.log(chalk.dim(`Removed old AIRchetipo configuration from Qwen Code`));
    }

    if (await fs.pathExists(oldAIRchetipoDir)) {
      await fs.remove(oldAIRchetipoDir);
      console.log(chalk.dim(`Removed old AIRCHETIPO configuration from Qwen Code`));
    }
  }
}

module.exports = { QwenSetup };
