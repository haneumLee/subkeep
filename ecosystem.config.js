/**
 * PM2 Ecosystem Configuration for SubKeep
 * 
 * Usage:
 *   Development: pm2 start ecosystem.config.js --env development
 *   Production:  pm2 start ecosystem.config.js --env production
 *   Staging:     pm2 start ecosystem.config.js --env staging
 * 
 * Commands:
 *   pm2 logs subkeep-backend    # View logs
 *   pm2 monit                    # Monitor processes
 *   pm2 restart all              # Restart all
 *   pm2 reload all               # Zero-downtime reload
 *   pm2 stop all                 # Stop all
 *   pm2 delete all               # Delete all
 */

module.exports = {
  apps: [
    // ==================== Backend API Server ====================
    {
      name: 'subkeep-backend',
      cwd: './backend',
      script: './bin/server',
      // Go binary doesn't need interpreter
      
      // Environment-specific variables
      env_development: {
        NODE_ENV: 'development',
        ENV: 'development',
        SERVER_PORT: 8080,
        LOG_LEVEL: 'debug',
      },
      env_staging: {
        NODE_ENV: 'staging',
        ENV: 'staging',
        SERVER_PORT: 8080,
        LOG_LEVEL: 'info',
      },
      env_production: {
        NODE_ENV: 'production',
        ENV: 'production',
        SERVER_PORT: 8080,
        LOG_LEVEL: 'warn',
      },
      
      // Process management
      instances: 2,              // Number of instances
      exec_mode: 'cluster',      // Cluster mode for load balancing
      autorestart: true,         // Auto restart on crash
      watch: false,              // Disable watch in production
      max_memory_restart: '500M', // Restart if memory exceeds 500MB
      
      // Logging
      error_file: './logs/backend-error.log',
      out_file: './logs/backend-out.log',
      log_date_format: 'YYYY-MM-DD HH:mm:ss Z',
      merge_logs: true,
      
      // Advanced
      min_uptime: '10s',         // Min uptime before considered online
      max_restarts: 10,          // Max restarts within 1 minute
      restart_delay: 4000,       // Delay between restarts
      kill_timeout: 5000,        // Time to wait for graceful shutdown
      listen_timeout: 3000,      // Time to wait for app to listen
      
      // Health monitoring
      vizion: true,              // Enable version control metadata
      post_update: ['echo "Backend updated"'],
      
      // Error handling
      exp_backoff_restart_delay: 100,
    },
    
    // ==================== Frontend (if using SSR) ====================
    // Uncomment if using Next.js SSR or similar
    /*
    {
      name: 'subkeep-frontend',
      cwd: './frontend',
      script: 'npm',
      args: 'start',
      
      env_development: {
        NODE_ENV: 'development',
        PORT: 3000,
      },
      env_production: {
        NODE_ENV: 'production',
        PORT: 3000,
      },
      
      instances: 1,
      exec_mode: 'cluster',
      autorestart: true,
      watch: false,
      max_memory_restart: '300M',
      
      error_file: './logs/frontend-error.log',
      out_file: './logs/frontend-out.log',
      log_date_format: 'YYYY-MM-DD HH:mm:ss Z',
    },
    */
  ],
  
  // ==================== Deployment Configuration ====================
  deploy: {
    // Production deployment to OracleVM
    production: {
      user: 'deploy',                        // SSH user
      host: ['192.168.1.100'],               // OracleVM server IP (change this)
      // host: ['subkeep.yourdomain.com'],   // Or use domain
      ref: 'origin/main',                    // Git branch
      repo: 'git@github.com:yourusername/subkeep.git', // Git repository
      path: '/home/deploy/subkeep',          // Deployment path
      ssh_options: 'StrictHostKeyChecking=no',
      
      // Pre-deployment commands
      'pre-deploy': [
        'git fetch --all',
        'git reset --hard origin/main',
      ].join(' && '),
      
      // Post-deployment commands
      'post-deploy': [
        'cd backend',
        'go build -o bin/server ./cmd/server',
        'pm2 reload ecosystem.config.js --env production',
        'pm2 save',
      ].join(' && '),
      
      // Pre-setup (first time only)
      'pre-setup': [
        'mkdir -p /home/deploy/subkeep',
        'mkdir -p /home/deploy/subkeep/logs',
      ].join(' && '),
    },
    
    // Staging deployment
    staging: {
      user: 'deploy',
      host: ['192.168.1.101'],               // Staging server IP
      ref: 'origin/dev',
      repo: 'git@github.com:yourusername/subkeep.git',
      path: '/home/deploy/subkeep-staging',
      ssh_options: 'StrictHostKeyChecking=no',
      
      'post-deploy': [
        'cd backend',
        'go build -o bin/server ./cmd/server',
        'pm2 reload ecosystem.config.js --env staging',
      ].join(' && '),
    },
  },
};
