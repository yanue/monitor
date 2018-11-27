FROM centos

# set repo
RUN sed -i "s|plugins=1|plugins=0|g" /etc/yum.conf\
    && mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup\
	&& curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo\
	&& rpm -ivh https://mirrors.aliyun.com/centos/7.5.1804/extras/x86_64/Packages/epel-release-7-11.noarch.rpm\
	&& yum makecache

# install base
RUN yum install -y net-tools vim wget which telnet killall pstree htop psmisc zsh git autojump\
	&& yum groupinstall -y "development tools"

# install zsh
RUN cd ~ && curl -O https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh\
    && chmod +x ./install.sh\
    && ./install.sh

# zsh plugins
RUN rm -fr ~/.oh-my-zsh/plugins/git\
    && git clone git://github.com/zsh-users/zsh-autosuggestions.git $ZSH_CUSTOM/plugins/zsh-autosuggestions\
    && git clone https://github.com/zsh-users/zsh-syntax-highlighting.git $ZSH_CUSTOM/plugins/zsh-syntax-highlighting\
    && echo "source $ZSH_CUSTOM/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh" >> ~/.zshrc\
    && echo "source $ZSH_CUSTOM/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh" >> ~/.zshrc\
    && echo "source /etc/profile.d/autojump.sh" >> ~/.zshrc

# set env
RUN mkdir /opt/go/bin -p\
    && mkdir /opt/go/pkg -p\
    && echo 'alias vi="vim"' >> /etc/zshrc\
	&& echo 'alias nets="netstat -ntpl"' >> /etc/zshrc\
	&& echo 'alias g="cd /opt/go/src"' >> /etc/zshrc\
	&& echo "PATH=/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin" >> /etc/zshrc\
	&& echo 'export GOPATH=/opt/go' >> /etc/zshrc\
	&& echo 'export GOROOT=/usr/local/go' >> /etc/zshrc\
	&& echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> /etc/zshrc

# install go
RUN cd /usr/local/src\
	&& curl -O https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz\
 	&& tar -xzf go1.11.2.linux-amd64.tar.gz\
 	&& mv go /usr/local/go\
 	&& ln -sf /usr/local/go/bin/go /usr/bin/go

# install gowatch
CMD exec zsh && source /etc/zshrc && go get github.com/silenceper/gowatch

# export http_proxy=http://127.0.0.1:1087;export https_proxy=http://127.0.0.1:1087;
# docker build -t centos-go .
# docker run -d --privileged -v /Users/yanue/golang/src:/opt/go/src --name test --hostname=test -p 3812:3812 centos-go /usr/sbin/init
