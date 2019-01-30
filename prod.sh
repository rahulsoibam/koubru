export API_PORT=':80'
# PostgreSQL related env variables
# DB_HOST='koubru-dev-db.caawnufl2ajb.ap-south-1.rds.amazonaws.com'
export DB_HOST='koubru-dev-db-cluster.cluster-caawnufl2ajb.ap-south-1.rds.amazonaws.com'
# DB_USERNAME='rahulsoibam'
export DB_USERNAME='rahulsoibam'
export DB_PASSWORD='058627439'
export DB_NAME='koubru_prod'
export DB_PORT='5432'

# User credentials postgresql database for security
export AUTH_DB_HOST='koubru-dev-db-cluster.cluster-caawnufl2ajb.ap-south-1.rds.amazonaws.com'
export AUTH_DB_USERNAME='rahulsoibam'
export AUTH_DB_PASSWORD='058627439'
export AUTH_DB_NAME='auth'
export AUTH_DB_PORT='5432'

# Credentials for redis instance for authentication
export AUTH_REDIS_ADDRESS='koubruauthredis-0001-001.qwpnzg.0001.aps1.cache.amazonaws.com:6379'
export AUTH_REDIS_PASSWORD=''

export CACHE_REDIS_ADDRESS='koubruauthredis-0001-002.qwpnzg.0001.aps1.cache.amazonaws.com'
export CACHE_REDIS_PASSWORD=''
# AWS related env variables
export AWS_REGION='ap-south-1'
export S3_BUCKET='media.rahulsoibam.me'

# SendGrid Key
export SENDGRID_API_KEY='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'

# Facebook Credentials
export FB_CLIENT_SECRET='fc5879613551a3ffd06a65ac0bb1f7b4'