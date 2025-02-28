- name: Setup EventBridge Monitoring for Allegro DB Connections
  hosts: localhost
  gather_facts: no
  vars:
    db_instance_id: "your-allegro-db-instance"
    event_rule_name: "AllegroDB-Connection-Monitor"
    sns_topic_name: "AllegroDBAlerts"
    sns_display_name: "Allegro DB Alerts"
    sns_email: "your-email@example.com"
    aws_region: "us-east-1"

  tasks:
    - name: Ensure AWS collections are installed
      ansible.builtin.command: ansible-galaxy collection install amazon.aws
      changed_when: false

    # Step 1: Enable CloudWatch Logs for Allegro DB
    - name: Enable CloudWatch Logs for Allegro DB
      amazon.aws.rds_instance:
        region: "{{ aws_region }}"
        db_instance_identifier: "{{ db_instance_id }}"
        enabled_cloudwatch_logs_exports:
          - "audit"
          - "error"
          - "general"
          - "slowquery"
      register: rds_logs

    - debug:
        msg: "CloudWatch Logs enabled for {{ db_instance_id }}"

    # Step 2: Create an SNS Topic for Alerts
    - name: Create SNS Topic
      amazon.aws.sns_topic:
        name: "{{ sns_topic_name }}"
        display_name: "{{ sns_display_name }}"
        region: "{{ aws_region }}"
      register: sns_topic

    - name: Subscribe Email to SNS Topic
      amazon.aws.sns_topic_subscription:
        topic: "{{ sns_topic_name }}"
        protocol: "email"
        endpoint: "{{ sns_email }}"
        region: "{{ aws_region }}"

    # Step 3: Create EventBridge Rule
    - name: Create EventBridge Rule for Allegro DB Connections
      amazon.aws.eventbridge_rule:
        name: "{{ event_rule_name }}"
        event_pattern: |
          {
            "source": ["aws.logs"],
            "detail-type": ["AWS API Call via CloudTrail"],
            "detail": {
              "eventSource": ["rds.amazonaws.com"],
              "eventName": ["Connect", "AuthorizeSecurityGroupIngress"],
              "requestParameters": {
                "dbInstanceIdentifier": ["{{ db_instance_id }}"]
              }
            }
          }
        state: "enabled"
        region: "{{ aws_region }}"

    # Step 4: Add SNS Target to EventBridge Rule
    - name: Add SNS Target to EventBridge Rule
      amazon.aws.eventbridge_target:
        rule: "{{ event_rule_name }}"
        arn: "{{ sns_topic.response.arn }}"
        region: "{{ aws_region }}"
      when: sns_topic is defined

    - debug:
        msg: "EventBridge Rule '{{ event_rule_name }}' created and linked to SNS"

