pipeline {
  agent any

  environment {
    DOCKER_REGISTRY = "${env.REGISTRY_URL ?: 'localhost:5000'}"
    SERVICE_NAME = 'api-gateway'
    GIT_BRANCH_NAME = "${env.BRANCH_NAME}"
  }

  stages {
    stage('Setup') {
      steps {
        script {
          // Auto-generate version tag from branch and build number
          def branchName = env.BRANCH_NAME.replaceAll('/', '-')
          def timestamp = new Date().format('yyyyMMdd-HHmmss')

          if (branchName == 'main') {
            env.BUILD_TAG = "v-${env.BUILD_NUMBER}-prod"
            env.ENVIRONMENT = 'production'
          } else if (branchName == 'develop' || branchName == 'dev') {
            env.BUILD_TAG = "v-${env.BUILD_NUMBER}-dev"
            env.ENVIRONMENT = 'dev'
          } else {
            env.BUILD_TAG = "v-${env.BUILD_NUMBER}-${branchName}-${timestamp}"
            env.ENVIRONMENT = 'dev'
          }

          echo "Building ${env.SERVICE_NAME} from branch ${branchName}"
          echo "Tag: ${env.BUILD_TAG}"
          echo "Environment: ${env.ENVIRONMENT}"
        }
      }
    }

    stage('Build Docker Image') {
      steps {
        script {
          sh """
            docker build -t ${env.DOCKER_REGISTRY}/${env.SERVICE_NAME}:${env.BUILD_TAG} .
            docker tag ${env.DOCKER_REGISTRY}/${env.SERVICE_NAME}:${env.BUILD_TAG} \
                       ${env.DOCKER_REGISTRY}/${env.SERVICE_NAME}:latest-${env.ENVIRONMENT}
            docker tag ${env.DOCKER_REGISTRY}/${env.SERVICE_NAME}:${env.BUILD_TAG} \
                       ${env.DOCKER_REGISTRY}/${env.SERVICE_NAME}:${env.GIT_BRANCH_NAME}-latest
          """
        }
      }
    }

    stage('Push to Registry') {
      steps {
        script {
          sh """
            docker push ${env.DOCKER_REGISTRY}/${env.SERVICE_NAME}:${env.BUILD_TAG}
            docker push ${env.DOCKER_REGISTRY}/${env.SERVICE_NAME}:latest-${env.ENVIRONMENT}
            docker push ${env.DOCKER_REGISTRY}/${env.SERVICE_NAME}:${env.GIT_BRANCH_NAME}-latest
          """
        }
      }
    }

    stage('Tag Git Commit') {
      steps {
        script {
          sh """
            git tag -a ${env.BUILD_TAG} -m "Build ${env.SERVICE_NAME} ${env.BUILD_TAG} for ${env.ENVIRONMENT}" || echo "Tag already exists"
            git push origin ${env.BUILD_TAG} || echo "Push failed or tag exists"
          """
        }
      }
    }
  }

  post {
    success {
      echo "✅ Successfully built ${env.SERVICE_NAME}:${env.BUILD_TAG}"
      echo "Image: ${env.DOCKER_REGISTRY}/${env.SERVICE_NAME}:${env.BUILD_TAG}"
    }
    failure {
      echo "❌ Build failed for ${env.SERVICE_NAME}:${env.BUILD_TAG}"
    }
    always {
      sh "docker image prune -f || true"
    }
  }
}
