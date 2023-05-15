package service

const (
    //instances
    MethodRunInstance            = "RUN_INSTANCE"
    MethodTerminateInstance      = "TERMINATE_INSTANCE"
    MethodDescribeInstance       = "DESCRIBE_INSTANCE"
    MethodTerminationProtection  = "MANAGE_TERMINATION_PROTECTION"
    MethodGetPlacementParameters = "ADDITIONAL_PARAM_ACTION"

    //images
    MethodCreateImage   = "CREATE_IMAGE"
    MethodDeleteImage   = "DELETE_IMAGE"
    MethodDescribeImage = "DESCRIBE_IMAGE"

    //volumes
    MethodCreateVolume          = "CREATE_VOLUME"
    MethodCreateAndAttachVolume = "CREATE_AND_ATTACH_VOLUME"
    MethodDeleteVolume          = "REMOVE_VOLUME"
    MethodDescribeVolume        = "DESCRIBE_VOLUME"

    //scripts
    MethodCreateScript   = "UPLOAD_SCRIPT"
    MethodDeleteScript   = "REMOVE_SCRIPT"
    MethodDescribeScript = "DESCRIBE_SCRIPT"

    //schedules
    MethodCreateSchedule   = "CREATE_SCHEDULE"
    MethodDescribeSchedule = "DESCRIBE_SCHEDULES"
    MethodDeleteSchedule   = "DELETE_SCHEDULE"

    //keypair
    MethodCreateKeypair   = "ADD_KEY"
    MethodDeleteKeypair   = "DELETE_KEY"
    MethodDescribeKeypair = "DESCRIBE_KEYS"

    //tags
    MethodUpdateTags = "UPDATE_TAGS"
    MethodDeleteTags = "DELETE_TAGS"

    //chef
    MethodGetChefProfiles = "GET_DEFAULT_REGION_CHEF_PROFILES"
)
